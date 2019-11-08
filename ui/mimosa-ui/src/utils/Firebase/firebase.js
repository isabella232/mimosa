import firebase from 'firebase';
import app from 'firebase/app';
import 'firebase/auth';

const conf = {
  apiKey: "AIzaSyCQieKOS6B36ut_o5n0loeW8rXetEqXnb0",
  authDomain: "mimosa-256008.firebaseapp.com",
  databaseURL: "https://mimosa-256008.firebaseio.com",
  projectId: "mimosa-256008",
  storageBucket: "mimosa-256008.appspot.com",
  messagingSenderId: "126377560493",
  appId: "1:126377560493:web:fb692f6332abe8a4bdb924"
};

class FirebaseApp {
  constructor() {
    app.initializeApp(conf);

    this.app = app;
    this.auth = app.auth();
    this.googleProv =  new firebase.auth.GoogleAuthProvider();
    this.db = app.firestore();
  }
}

export default FirebaseApp;