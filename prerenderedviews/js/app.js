import React from "react";
import ReactDOM from "react-dom";
import { render } from "react-dom";
import { BrowserRouter, Route } from 'react-router-dom'
import {browserHistory} from 'react-router';
import "./../css/index.css"
import $ from 'jquery';
import "bulma"
import auth0 from 'auth0-js';
const AUTH0_CLIENT_ID = "Y4lTZL7LZ05OnNglAcsmogfmTbDPDbDN";
const AUTH0_DOMAIN = "cryptex.auth0.com";
const AUTH0_CALLBACK_URL = location.href;
const AUTH0_API_AUDIENCE = "https://cryptex.auth0.com/api/v2/";


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
	      if (
	        authResult !== null &&
	        authResult.accessToken !== null &&
	        authResult.idToken !== null
	      ) {
	        localStorage.setItem("access_token", authResult.accessToken);
	        localStorage.setItem("id_token", authResult.idToken);
	        localStorage.setItem(
	          "profile",
	          JSON.stringify(authResult.idTokenPayload)
	        );
	        window.location = window.location.href.substr(
	          0,
	          window.location.href.indexOf("#")
	        );
	      }
	    });
	  }

	  setup() {
	    $.ajaxSetup({
	      beforeSend: (r) => {
	        if (localStorage.getItem("access_token")) {
	          r.setRequestHeader(
	            "Authorization",
	            "Bearer " + localStorage.getItem("access_token")
	          );
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
			return (<LoggedIn />);
		else
			return (<Home />);
	}
}
class LoggedIn extends React.Component {
	render() {
		return(<p>You are logged in</p>);
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
      <div className="container">
        <div className="row">
          <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
            <h1>Jokeish</h1>
            <p>A load of Dad jokes XD</p>
            <p>Sign in to get access </p>
            <a
              onClick={this.authenticate}
              className="btn btn-primary btn-lg btn-login btn-block"
            >
              Sign In
            </a>
          </div>
        </div>
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
	<BrowserRouter>
		<div>
		<Route path="/" component={App} history={browserHistory}/>
		<Route path="/callback" component={Callback} history={browserHistory}/>
		</div>
	</BrowserRouter>
,document.getElementById('app'));


if (module.hot)
{
	module.hot.accept();
}