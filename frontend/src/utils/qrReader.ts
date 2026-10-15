import jsQR from "jsqr";

/**
 * Extracts a clean token/UUID from decoded QR string, handling URLs, query params, JSON, or raw tokens.
 */
export function extractCheckInToken(rawText: string): string {
  if (!rawText) return "";
  let text = rawText.trim();

  // 1. Try URL parsing
  try {
    const url = new URL(text);
    const candidate =
      url.searchParams.get("checkin") ||
      url.searchParams.get("token") ||
      url.searchParams.get("id") ||
      url.searchParams.get("qr_token");
    if (candidate) {
      return candidate.trim();
    }
  } catch {
    // Not a valid URL, check relative query string
    if (text.includes("?")) {
      try {
        const queryPart = text.split("?")[1];
        const params = new URLSearchParams(queryPart);
        const candidate =
          params.get("checkin") ||
          params.get("token") ||
          params.get("id") ||
          params.get("qr_token");
        if (candidate) {
          return candidate.trim();
        }
      } catch {
        // ignore
      }
    }
  }

  // 2. Try JSON parsing
  if (text.startsWith("{") && text.endsWith("}")) {
    try {
      const parsed = JSON.parse(text);
      if (parsed.qr_token) return String(parsed.qr_token).trim();
      if (parsed.token) return String(parsed.token).trim();
      if (parsed.id) return String(parsed.id).trim();
      if (parsed.checkin) return String(parsed.checkin).trim();
    } catch {
      // ignore
    }
  }

  return text;
}

/**
 * Decodes a QR code from any image File or Blob with high-precision multi-pass analysis.
 * Handles full screenshots, cropped tickets, and colored/beige backgrounds seamlessly.
 */
export async function decodeQrFromImageFile(file: File | Blob): Promise<string | null> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    img.onload = () => {
      URL.revokeObjectURL(objectUrl);
      try {
        const naturalWidth = img.naturalWidth || img.width;
        const naturalHeight = img.naturalHeight || img.height;

        if (naturalWidth === 0 || naturalHeight === 0) {
          resolve(null);
          return;
        }

        const canvas = document.createElement("canvas");
        canvas.width = naturalWidth;
        canvas.height = naturalHeight;
        const ctx = canvas.getContext("2d", { willReadFrequently: true });
        if (!ctx) {
          resolve(null);
          return;
        }

        ctx.drawImage(img, 0, 0, naturalWidth, naturalHeight);

        // --- PASS 1: Full Image with jsQR ---
        let imageData = ctx.getImageData(0, 0, naturalWidth, naturalHeight);
        let code = jsQR(imageData.data, naturalWidth, naturalHeight, {
          inversionAttempts: "dontInvert",
        });
        if (code && code.data) {
          resolve(code.data);
          return;
        }

        // --- PASS 2: Full Image with Inversion (Dark/Light) ---
        code = jsQR(imageData.data, naturalWidth, naturalHeight, {
          inversionAttempts: "attemptBoth",
        });
        if (code && code.data) {
          resolve(code.data);
          return;
        }

        // --- PASS 3: Ticket Center Region Crop ---
        // In e-tickets, the QR code is centered vertically around 15% - 75%
        const cropX = Math.floor(naturalWidth * 0.15);
        const cropY = Math.floor(naturalHeight * 0.15);
        const cropW = Math.floor(naturalWidth * 0.70);
        const cropH = Math.floor(naturalHeight * 0.65);

        const cropCanvas = document.createElement("canvas");
        cropCanvas.width = cropW;
        cropCanvas.height = cropH;
        const cropCtx = cropCanvas.getContext("2d", { willReadFrequently: true });

        if (cropCtx) {
          cropCtx.drawImage(img, cropX, cropY, cropW, cropH, 0, 0, cropW, cropH);
          const cropData = cropCtx.getImageData(0, 0, cropW, cropH);

          code = jsQR(cropData.data, cropW, cropH, {
            inversionAttempts: "attemptBoth",
          });
          if (code && code.data) {
            resolve(code.data);
            return;
          }

          // --- PASS 4: Contrast & Binarization on Cropped Area ---
          // Boost contrast to remove beige/colored ticket backgrounds
          const data = cropData.data;
          for (let i = 0; i < data.length; i += 4) {
            const gray = data[i] * 0.299 + data[i + 1] * 0.587 + data[i + 2] * 0.114;
            // Threshold at 150 to cleanly separate dark QR modules from light/beige background
            const val = gray < 150 ? 0 : 255;
            data[i] = val;
            data[i + 1] = val;
            data[i + 2] = val;
          }
          cropCtx.putImageData(cropData, 0, 0);

          code = jsQR(cropData.data, cropW, cropH, {
            inversionAttempts: "attemptBoth",
          });
          if (code && code.data) {
            resolve(code.data);
            return;
          }
        }

        // --- PASS 5: Normalized Resolution (800px scale) ---
        if (naturalWidth > 1200 || naturalHeight > 1200) {
          const scale = 800 / Math.max(naturalWidth, naturalHeight);
          const scaledW = Math.floor(naturalWidth * scale);
          const scaledH = Math.floor(naturalHeight * scale);

          const scaleCanvas = document.createElement("canvas");
          scaleCanvas.width = scaledW;
          scaleCanvas.height = scaledH;
          const scaleCtx = scaleCanvas.getContext("2d", { willReadFrequently: true });
          if (scaleCtx) {
            scaleCtx.drawImage(img, 0, 0, scaledW, scaledH);
            const scaleData = scaleCtx.getImageData(0, 0, scaledW, scaledH);
            code = jsQR(scaleData.data, scaledW, scaledH, {
              inversionAttempts: "attemptBoth",
            });
            if (code && code.data) {
              resolve(code.data);
              return;
            }
          }
        }

        resolve(null);
      } catch (err) {
        console.error("Error in decodeQrFromImageFile:", err);
        resolve(null);
      }
    };

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl);
      resolve(null);
    };

    img.src = objectUrl;
  });
}
