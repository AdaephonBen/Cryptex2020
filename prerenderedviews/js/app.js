import React from "react";
import ReactDOM from "react-dom";
import { render } from "react-dom";
import { BrowserRouter, Route } from 'react-router-dom'
import {browserHistory} from 'react-router';
import "./../css/index.css"

import $ from 'jquery';
import auth0 from 'auth0-js';

const AUTH0_CLIENT_ID = "xSWF7EZ8NNiusQpwCeKbh21TGjRR7tIy";
const AUTH0_DOMAIN = "cryptex2020.auth0.com";
const AUTH0_CALLBACK_URL = "http://localhost:8080";
const AUTH0_API_AUDIENCE = "https://cryptex2020.auth0.com/api/v2/";

class App extends React.Component {
	parseHash() {
	    this.auth0 = new auth0.WebAuth({
	      domain: AUTH0_DOMAIN,
	      clientID: AUTH0_CLIENT_ID
	    });
	    this.auth0.parseHash(window.location.hash, (err, authResult) => {
	      if (err) {
	        return console.log(err);
	      }
	      if(authResult !== null && authResult.accessToken !== null && authResult.idToken !== null){
	              localStorage.setItem('access_token', authResult.accessToken);
	              localStorage.setItem('id_token', authResult.idToken);
	              localStorage.setItem('profile', JSON.stringify(authResult.idTokenPayload));
	          window.location = window.location.href.substr(0, window.location.href.indexOf(''))
	            }
	    });
	  }

	  setup() {
	    $.ajaxSetup({
	          'beforeSend': function(xhr) {
	            if (localStorage.getItem('access_token')) {
	              xhr.setRequestHeader('Authorization',
	                    'Bearer ' + localStorage.getItem('access_token'));
	            }
	          }
	        });
	  }

	  setState() {
	    let idToken = localStorage.getItem("id_token");
	    if (idToken) {
	      this.loggedIn = true;
	    } else {
	      this.loggedIn = false;
	    }
	  }

	  componentWillMount() {
	    this.setup();
	    this.parseHash();
	    this.setState();
	  }

	render() {
		if (this.loggedIn)
			return (<div><nav>
    		<a href="/rules/" onClick=""><div className="nav-links">Rules</div></a>
    		<a href="/crypt2019/"><div className="nav-links">Cryptex 2019</div></a>
    		<a href="/sponsors"><div className="nav-links">Sponsors</div></a>
    		<a href="/about/"><div className="nav-links">About Us</div></a>
    	</nav> <LoggedIn /> </div> );
		else
			return (<div><nav>
    		<a href="/rules/" onClick=""><div className="nav-links">Rules</div></a>
    		<a href="/crypt2019/"><div className="nav-links">Cryptex 2019</div></a>
			<a href="/" onClick=""><div className="nav-links">C R Y P T E X</div></a>
    		<a href="/sponsors"><div className="nav-links">Sponsors</div></a>
    		<a href="/about/"><div className="nav-links">About Us</div></a>
    	</nav><Home /></div>);
	}
}
class LoggedIn extends React.Component 
{
	render() 
	{
		return(<div><p>You are logged in, {JSON.parse(localStorage.getItem("profile")).name}. </p><p>Give us a username</p></div>);
	}
}
class Home extends React.Component {
  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }
  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid profile",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL
    });
    this.WebAuth.authorize();
  }

  render() {
    return (
    	<div>
    	<button className="DiveInButton" onClick={this.authenticate}><div className="transform">D I V E &nbsp; I N</div></button>
    	</div>
    	);
  }
}

class Callback extends React.Component {
	render() {
		return(<h1>Loading</h1>);
	}
}

render(
	<App />
,document.getElementById('app'));


if (module.hot)
{
	module.hot.accept();
}