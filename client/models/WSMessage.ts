interface WSMessage {
  type: string;
  message: string;
  error: Error;
  submissionId: number;
}

export default WSMessage;
