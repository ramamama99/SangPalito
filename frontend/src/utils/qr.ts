import QRCode from "qrcode";

export async function generateQrDataUrl(ticketIdOrToken: string): Promise<string> {
  if (!ticketIdOrToken) return "";
  
  // Format checkin URL pointing to /scan so scanning opens the real check-in scanner
  const checkinUrl = ticketIdOrToken.startsWith("http")
    ? ticketIdOrToken
    : `${window.location.origin}/scan?checkin=${encodeURIComponent(ticketIdOrToken)}`;

  return QRCode.toDataURL(checkinUrl, {
    width: 280,
    margin: 2,
    color: {
      dark: "#0f172a",
      light: "#ffffff",
    },
    errorCorrectionLevel: "M",
  });
}
