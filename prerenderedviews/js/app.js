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

export default function Level({ clientID }) {
	return (
		<Query query={GET_LEVEL} variables = {{ clientID }}>
			{({ data, loading, error }) => {
				if (loading)	return <Loading />;
				if (error)	return (<p>Error : {error.message} </p>)
			}}
		</Query>
	);
}

// to do
// immortal db for emails

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
	              localStorage.setItem('email', JSON.stringify(authResult.idTokenPayload));
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
	 renderBody()
	 {
		if (this.loggedIn)
			return (<div><nav>
    		<a href="/rules/" onClick=""><div className="nav-links">Rules</div></a>
    		<a href="/crypt2019/"><div className="nav-links">Cryptex 2019</div></a>
			<a href="/" onClick=""><div id="nav-links-main">C R Y P T E X</div></a>
    		<a href="/sponsors"><div className="nav-links">Sponsors</div></a>
    		<a href="/about/"><div className="nav-links">About Us</div></a>
    		<div className = "bar"></div>
    	</nav> <LoggedIn /> </div> );
		else
			return (
			<div>	
				<nav>
		    		<a href="/rules/" onClick=""><div className="nav-links">Rules</div></a>
		    		<a href="/crypt2019/"><div className="nav-links">Cryptex 2019</div></a>
					<a href="/" onClick=""><div id="nav-links-main">C R Y P T E X</div></a>
		    		<a href="/sponsors"><div className="nav-links">Sponsors</div></a>
		    		<a href="/about/"><div className="nav-links">About Us</div></a>
    			</nav>
    			<Home />
			</div>
    	);
	 }
	render() {
		return this.loggedIn==undefined ? (<div className="loader"></div>) : this.renderBody() ;
	}
}
class LoggedIn extends React.Component 
{
	constructor(props)
	{
		super(props);
		this.state={value: "", level:"", client:{}};
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
		this.fetchLevel = this.fetchLevel.bind(this);
	}
	handleChange(event)
	{
		this.setState({value : event.target.value});
	}
	handleSubmit(event)
	{
		event.preventDefault();
		let url = "http://localhost:8080/graphql?query={doesUsernameExist(username:\"" + this.state.value + "\")}";
		fetch(url).then(response => response.json())
		.then(result => {
			if (result.data.doesUsernameExist == true)
			{
				alert("That username exists");
			}
			else
			{
				var loginUrl = "/adduser/"+JSON.parse(localStorage.getItem("email")).email+"/"+this.state.value+"/"+localStorage.getItem("id_token");
				fetch(loginUrl).then(() => {
					let user = new User(this.state.level, localStorage.getItem("id_token"), this.state.value);
					this.setState({client: user});
					let url = "http://localhost:8080/graphql?query={level(clientID:\"" + JSON.parse(localStorage.getItem("email")).email + "\")}"
					this.fetchLevel();
				});
			}
		});
	}
	fetchLevel()
	{
		let url = "http://localhost:8080/graphql?query={level(clientID:\"" + JSON.parse(localStorage.getItem("email")).email + "\")}"
		fetch(url)
		.then(response => response.json())
		.then(result => {
			this.setState({level: result.data.level});
		});
	}
	componentDidMount() {
		this.fetchLevel();
	}
	renderUsername()
	{
		return(
			<div class="username-form">
				<p>You are logged in, {JSON.parse(localStorage.getItem("email")).email}. </p>
				<p>Give us a username.</p>
				<form onSubmit={this.handleSubmit}>
					<input type="name" class="username" value={this.state.value} onChange={this.handleChange}/>
					<br /><br />
					<input type="submit" class="username-button" value="Submit" />
				</form>
			</div>
		);
	}
	renderRules()
	{
		return(
			<div className="rules-container">
				<h1 className = "rules">Rules Shit Here</h1>
				<button className="accept-rules">I accept all this shi</button>
			</div>
		);
	}
	render() 
	{
		const level = this.state.level ;
		if (level)
		{
			switch(level)
			{
				case "-2":
					return(this.renderUsername());
					break;
				case "-1":
					return(this.renderRules());
					break;
			}
		}
		else
		{
			return(<div className="loader"></div>);
		}
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
      scope: "openid email",
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
class User {
	constructor(level, clientSecret, username)
	{
		this.level = level ;
		this.clientSecret = clientSecret ;
		this.username = username ;
	}
	getLevel()
	{
		return(this.level);
	}
	getClientSecret()
	{
		return(this.clientSecret);
	}
	getUsername()
	{
		return(this.username);
	}
}

render(
	<App />
,document.getElementById('app'));


if (module.hot)
{
	module.hot.accept();
}