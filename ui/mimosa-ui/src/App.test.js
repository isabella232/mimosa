import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import FirebaseApp, { FirebaseContext } from './utils/Firebase';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <FirebaseContext.Provider value={new FirebaseApp()}>
      <App />
    </FirebaseContext.Provider>
    , div);
  ReactDOM.unmountComponentAtNode(div);
});
