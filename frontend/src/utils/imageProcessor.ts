export interface PreprocessOptions {
  mode?: 'adaptive' | 'binarize' | 'contrast' | 'original'
  contrast?: number // -100 to 100
  sharpen?: boolean
  rotation?: number // 0, 90, 180, 270
  targetWidth?: number // default 1400
}

/**
 * Memutar canvas berdasarkan derajat (0, 90, 180, 270)
 */
export function rotateCanvas(sourceCanvas: HTMLCanvasElement, degrees: number): HTMLCanvasElement {
  const normalizedDeg = ((degrees % 360) + 360) % 360
  if (normalizedDeg === 0) return sourceCanvas

  const rotated = document.createElement('canvas')
  const ctx = rotated.getContext('2d')
  if (!ctx) return sourceCanvas

  if (normalizedDeg === 90 || normalizedDeg === 270) {
    rotated.width = sourceCanvas.height
    rotated.height = sourceCanvas.width
  } else {
    rotated.width = sourceCanvas.width
    rotated.height = sourceCanvas.height
  }

  ctx.translate(rotated.width / 2, rotated.height / 2)
  ctx.rotate((normalizedDeg * Math.PI) / 180)
  ctx.drawImage(sourceCanvas, -sourceCanvas.width / 2, -sourceCanvas.height / 2)

  return rotated
}

/**
 * Memotong area tertentu dari Video atau Canvas (misal: area viewfinder KTP)
 */
export function cropSourceToCanvas(
  source: HTMLVideoElement | HTMLCanvasElement | HTMLImageElement,
  cropRect: { xPercent: number; yPercent: number; widthPercent: number; heightPercent: number }
): HTMLCanvasElement {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')

  const naturalWidth = 'videoWidth' in source ? source.videoWidth : ('naturalWidth' in source ? source.naturalWidth : source.width)
  const naturalHeight = 'videoHeight' in source ? source.videoHeight : ('naturalHeight' in source ? source.naturalHeight : source.height)

  const sx = Math.max(0, Math.round(cropRect.xPercent * naturalWidth))
  const sy = Math.max(0, Math.round(cropRect.yPercent * naturalHeight))
  const sWidth = Math.min(naturalWidth - sx, Math.round(cropRect.widthPercent * naturalWidth))
  const sHeight = Math.min(naturalHeight - sy, Math.round(cropRect.heightPercent * naturalHeight))

  canvas.width = sWidth
  canvas.height = sHeight

  if (ctx) {
    ctx.drawImage(source, sx, sy, sWidth, sHeight, 0, 0, sWidth, sHeight)
  }

  return canvas
}

/**
 * Menghitung threshold optimal menggunakan algoritma Otsu (Binarization)
 */
function getOtsuThreshold(grayData: Uint8Array): number {
  const histogram = new Array(256).fill(0)
  for (let i = 0; i < grayData.length; i++) {
    histogram[grayData[i]]++
  }

  const total = grayData.length
  let sum = 0
  for (let i = 0; i < 256; i++) {
    sum += i * histogram[i]
  }

  let sumB = 0
  let wB = 0
  let wF = 0
  let varMax = 0
  let threshold = 128

  for (let t = 0; t < 256; t++) {
    wB += histogram[t]
    if (wB === 0) continue
    wF = total - wB
    if (wF === 0) break

    sumB += t * histogram[t]
    const mB = sumB / wB
    const mF = (sum - sumB) / wF

    const varBetween = wB * wF * (mB - mF) * (mB - mF)
    if (varBetween > varMax) {
      varMax = varBetween
      threshold = t
    }
  }

  return threshold
}

/**
 * Filter penajaman (Sharpening convolution 3x3)
 */
function applySharpen(data: Uint8ClampedArray, width: number, height: number): void {
  // Salin data asli untuk referensi
  const copy = new Uint8ClampedArray(data)
  // Kernel penajaman halus:
  // [  0, -0.5,  0 ]
  // [ -0.5, 3.0, -0.5 ]
  // [  0, -0.5,  0 ]
  for (let y = 1; y < height - 1; y++) {
    for (let x = 1; x < width - 1; x++) {
      const idx = (y * width + x) * 4
      const top = ((y - 1) * width + x) * 4
      const bottom = ((y + 1) * width + x) * 4
      const left = (y * width + (x - 1)) * 4
      const right = (y * width + (x + 1)) * 4

      for (let c = 0; c < 3; c++) {
        const val = 3.0 * copy[idx + c] - 0.5 * (copy[top + c] + copy[bottom + c] + copy[left + c] + copy[right + c])
        data[idx + c] = Math.max(0, Math.min(255, val))
      }
    }
  }
}

/**
 * Pipeline Pre-processing khusus e-KTP Indonesia:
 * 1. Rescaling ke resolusi optimal (width ~ 1400px)
 * 2. Supresi latar belakang biru/cyan khas e-KTP menggunakan pembobotan kanal warna
 * 3. Histogram Stretching & Contrast Normalization
 * 4. Sharpening untuk mengatasi blur kamera HP
 * 5. Binarisasi (Otsu) opsional jika teks tipis/kurang tegas
 */
export async function processImageForOcr(
  imageInput: File | Blob | string | HTMLCanvasElement,
  options: PreprocessOptions = {}
): Promise<{ dataUrl: string; width: number; height: number }> {
  const {
    mode = 'adaptive',
    contrast = 50,
    sharpen = true,
    rotation = 0,
    targetWidth = 1400
  } = options

  return new Promise((resolve, reject) => {
    const applyPipeline = (sourceCanvas: HTMLCanvasElement) => {
      let canvas = sourceCanvas

      // 1. Tangani Rotasi jika ada
      if (rotation !== 0) {
        canvas = rotateCanvas(canvas, rotation)
      }

      // 2. Normalisasi ukuran (Rescaling)
      // Tesseract v5 bekerja paling akurat jika teks memiliki tinggi x-height 30-40px (lebar kartu ~1400px)
      if (canvas.width !== targetWidth) {
        const scale = targetWidth / canvas.width
        const scaledCanvas = document.createElement('canvas')
        scaledCanvas.width = targetWidth
        scaledCanvas.height = Math.round(canvas.height * scale)
        const sCtx = scaledCanvas.getContext('2d')
        if (sCtx) {
          sCtx.imageSmoothingEnabled = true
          sCtx.imageSmoothingQuality = 'high'
          sCtx.drawImage(canvas, 0, 0, scaledCanvas.width, scaledCanvas.height)
          canvas = scaledCanvas
        }
      }

      if (mode === 'original') {
        resolve({
          dataUrl: canvas.toDataURL('image/png'),
          width: canvas.width,
          height: canvas.height
        })
        return
      }

      const ctx = canvas.getContext('2d')
      if (!ctx) {
        resolve({
          dataUrl: canvas.toDataURL('image/png'),
          width: canvas.width,
          height: canvas.height
        })
        return
      }

      const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
      const data = imageData.data
      const totalPixels = canvas.width * canvas.height
      const grayBuffer = new Uint8Array(totalPixels)

      // 3. Konversi Warna & Supresi Background KTP
      // Catatan: Latar KTP bernuansa biru/cyan (R rendah, G sedang, B tinggi).
      // Teks KTP berwarna hitam legam (R, G, B rendah).
      // Mengurangi bobot kanal biru dan menaikkan kanal merah membuat latar KTP menjadi putih pudar!
      for (let i = 0, p = 0; i < data.length; i += 4, p++) {
        const r = data[i]
        const g = data[i + 1]
        const b = data[i + 2]

        let gray: number
        if (mode === 'adaptive') {
          // Bobot penekan latar biru: Red 45%, Green 45%, Blue 10%
          gray = 0.45 * r + 0.45 * g + 0.1 * b
        } else {
          // Standar Luminance Grayscale
          gray = 0.299 * r + 0.587 * g + 0.114 * b
        }

        grayBuffer[p] = Math.round(gray)
      }

      // 4. Mode Binarisasi (Otsu Thresholding)
      if (mode === 'binarize') {
        const threshold = getOtsuThreshold(grayBuffer)
        for (let i = 0, p = 0; i < data.length; i += 4, p++) {
          const val = grayBuffer[p] < threshold ? 0 : 255
          data[i] = val
          data[i + 1] = val
          data[i + 2] = val
        }
      } else {
        // 5. Adaptive Contrast Stretching (Histogram Min-Max Stretch)
        // Cari nilai min dan max untuk membentangkan dinamis range kontras
        let min = 255
        let max = 0
        for (let p = 0; p < totalPixels; p++) {
          const v = grayBuffer[p]
          if (v < min) min = v
          if (v > max) max = v
        }

        // Cegah divide by zero
        const range = max - min > 10 ? max - min : 255
        const contrastFactor = (259 * (contrast + 255)) / (255 * (259 - contrast))

        for (let i = 0, p = 0; i < data.length; i += 4, p++) {
          // Normalisasi stretch
          let val = ((grayBuffer[p] - min) / range) * 255
          // Tambahkan kontras tambahan
          val = contrastFactor * (val - 128) + 128
          val = Math.max(0, Math.min(255, val))

          data[i] = val
          data[i + 1] = val
          data[i + 2] = val
        }

        // 6. Sharpening filter jika diaktifkan
        if (sharpen) {
          applySharpen(data, canvas.width, canvas.height)
        }
      }

      ctx.putImageData(imageData, 0, 0)
      resolve({
        dataUrl: canvas.toDataURL('image/png'),
        width: canvas.width,
        height: canvas.height
      })
    }

    if (imageInput instanceof HTMLCanvasElement) {
      applyPipeline(imageInput)
    } else if (typeof imageInput === 'string') {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = img.naturalWidth || img.width
        canvas.height = img.naturalHeight || img.height
        const ctx = canvas.getContext('2d')
        if (ctx) ctx.drawImage(img, 0, 0)
        applyPipeline(canvas)
      }
      img.onerror = reject
      img.src = imageInput
    } else {
      const reader = new FileReader()
      reader.onload = (e) => {
        const img = new Image()
        img.onload = () => {
          const canvas = document.createElement('canvas')
          canvas.width = img.naturalWidth || img.width
          canvas.height = img.naturalHeight || img.height
          const ctx = canvas.getContext('2d')
          if (ctx) ctx.drawImage(img, 0, 0)
          applyPipeline(canvas)
        }
        img.onerror = reject
        img.src = e.target?.result as string
      }
      reader.onerror = reject
      reader.readAsDataURL(imageInput)
    }
  })
}
