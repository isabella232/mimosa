import firebase from 'firebase';

const conf = {
  apiKey: "AIzaSyCQieKOS6B36ut_o5n0loeW8rXetEqXnb0",
  authDomain: "mimosa-256008.firebaseapp.com",
  databaseURL: "https://mimosa-256008.firebaseio.com",
  projectId: "mimosa-256008",
  storageBucket: "mimosa-256008.appspot.com",
  messagingSenderId: "126377560493",
  appId: "1:126377560493:web:fb692f6332abe8a4bdb924"
};

firebase.initializeApp(conf);

// export const provider = new firebase.auth.GoogleAuthProvider();
// provider.setCustomParameters({
//   hd: "puppet.com"
// });

export const googleProvider = new firebase.auth.GoogleAuthProvider();
export const auth = firebase.auth();
require("firebase/functions");
export const functions = firebase.functions();
export default firebase;