import QRCode from "qrcode";

export async function generateQrDataUrl(ticketId: string): Promise<string> {
  const checkinUrl = `${window.location.origin}${window.location.pathname}?checkin=${encodeURIComponent(ticketId)}`;
  return QRCode.toDataURL(checkinUrl, {
    width: 240,
    margin: 1,
    color: { dark: "#195A5F" },
  });
}
