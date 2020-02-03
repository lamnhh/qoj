interface User {
  username: string;
  fullname: string;
  profilePicture: string;
}

const emptyUser: User = {
  username: "",
  fullname: "",
  profilePicture: ""
};

export default User;
export { emptyUser };
