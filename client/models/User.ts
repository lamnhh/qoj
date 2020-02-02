interface User {
  username: string;
  fullname: string;
}

const emptyUser: User = {
  username: "",
  fullname: ""
};

export default User;
export { emptyUser };
