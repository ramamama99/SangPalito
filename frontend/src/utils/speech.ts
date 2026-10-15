export function speakQueueCall(queueNumber: string, poliName: string) {
  if (typeof window !== "undefined" && "speechSynthesis" in window) {
    window.speechSynthesis.cancel();
    const speech = new SpeechSynthesisUtterance(
      `Nomor antrean ${queueNumber.replace("-", " ")}, silakan menuju ${poliName}`,
    );
    speech.lang = "id-ID";
    window.speechSynthesis.speak(speech);
  }
}
