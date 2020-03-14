interface Contest {
  id: number;
  name: string;
  startDate: Date; // in ISO format, i.e. "2020-02-07T09:44:15.991Z"
  duration: number; // in minutes
  numberOfParticipants: number;
  isRegistered: boolean;
}

const defaultContest: Contest = {
  id: 0,
  name: "",
  startDate: new Date(),
  duration: 150,
  numberOfParticipants: 0,
  isRegistered: false
};

export default Contest;
export { defaultContest };
