import React from "react";
import "./App.css";

const App: React.FC = () => {
  return <div className="">Hello World</div>;
};

export default App;

// console.log(process.env.REACT_APP_NOT_SECRET_CODE);
// "proxy": "http://localhost:0000" in package.json
// import graphql from 'babel-plugin-relay/macro';

// useInterval hook to refresh token?
// https://overreacted.io/making-setinterval-declarative-with-react-hooks/

// const print = <T extends {}>(msg: T): T => {
//   return msg;
// };

// const print = function<T>(msg: T): T {
//   return msg;
// };

// console.log(print(84545));
