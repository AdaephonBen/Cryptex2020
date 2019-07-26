import React from "react";
import image from './../07.png'
import ReactDOM from "react-dom";
import { render } from "react-dom";
import "./../css/index.css"
import $ from 'jquery';
import auth0 from 'auth0-js';

const AUTH0_CLIENT_ID = "xSWF7EZ8NNiusQpwCeKbh21TGjRR7tIy";
const AUTH0_DOMAIN = "cryptex2020.auth0.com";
const AUTH0_CALLBACK_URL = "https://2020.elan.org.in";
const AUTH0_API_AUDIENCE = "https://cryptex2020.auth0.com/api/v2/";

export default function Level({ clientID }) {
	return (
		<Query query={GET_LEVEL} variables={{ clientID }}>
			{({ data, loading, error }) => {
				if (loading) return <Loading />;
				if (error) return (<p>Error : {error.message} </p>)
			}}
		</Query>
	);
}


// <div className="navbar" id="mainNavBar">
// 	<div className="container">
// 	<div className="row">
// 	<div className="two columns"><a><div className="nav-links">Rules</div></a></div>
//   		<div className="two columns"><a><div className="nav-links">Cryptex 2019</div></a></div>
// 	<div className="two columns"><a><div id="nav-links-main">C R Y P T E X</div></a></div>
//   		<div className="two columns"><a><div className="nav-links">Sponsors</div></a></div>
//   		<div className="three columns"><a><div className="nav-links">About Us</div></a></div>
//   		</div>
//   		</div>
// </div>
class Navbar extends React.Component {
	render() {
		return (
			<nav class="animated fadeInDown">
				<ul>
					<a href="/rules"><li>Guidelines</li></a>
					<a href="/" className="main animated flipInX"><li>C R Y P T E X</li></a>
					<a href="/leaderboardtable" id="leaderboard-nav-link"><li>Leaderboard</li></a>
				</ul>
			</nav>
		);
	}
}

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
			if (authResult !== null && authResult.accessToken !== null && authResult.idToken !== null) {
				localStorage.setItem('access_token', authResult.accessToken);
				localStorage.setItem('id_token', authResult.idToken);
				localStorage.setItem('email', JSON.stringify(authResult.idTokenPayload));
				window.location = window.location.href.substr(0, window.location.href.indexOf(''))

			}
		});
	}

	setup() {
		$.ajaxSetup({
			'beforeSend': function (xhr) {
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
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	componentWillMount() {
		this.setup();
		this.parseHash();
		this.setState();
	}
	renderBody() {
		if (this.loggedIn)
			return (<div> <Navbar /> <LoggedIn /></div>);
		else
			return (
				<div>
					<Navbar />
					<Home />
				</div>
			);
	}
	render() {
		return this.loggedIn == undefined ? (<div className="loader"></div>) : this.renderBody();
	}
}
class LoggedIn extends React.Component {
	constructor(props) {
		super(props);
		this.state = { value: "", level: "", client: {} };
		this.handleChange = this.handleChange.bind(this);
		this.fetchLevel = this.fetchLevel.bind(this);
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}

	fetchLevel() {
		let url = "https://2020.elan.org.in/whichlevel/"+JSON.parse(localStorage.getItem("email")).email ;
		fetch(url)
			.then(response => response.json())
			.then(result => {
				this.setState({ level: result.message });
			});
	}

	componentDidMount() {
		this.fetchLevel();
	}
	render() {
		const level = this.state.level;
		if (level) {
			switch (level) {
				case "-2":
					return (<LevelUsername />);
					break;
				case "-1":
					return (<LevelRules />);
					break;
				case "0":
					return (<LevelImage />);
				case "1":
					return (<LevelImage />);
				case "2":
					return (<LevelText />);
				case "3":
					return (<LevelMidi />);
				case "4":
					return (<LevelImage />);
				case "5":
					return (<LevelImage />);
				case "6":
					return (<LevelWon />);
			}
		}
		else {
			return (<div className="loader"></div>);
		}
	}
}
class LevelWon extends React.Component {
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	render() {
		return (
			<div className="won congrats">
				You have won. <br /> Congrats.
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
	}
}

class LevelImage extends React.Component {
	constructor(props) {
		super(props);
		this.state = { value: "", url: "", level: -3 };
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}
	handleSubmit(event) {
		event.preventDefault();
		let url = "https://2020.elan.org.in/answer/" + localStorage.getItem("id_token") + "/" + this.state.level.toString() + "/" + this.state.value;
		fetch(url).then(() => {
			window.location.reload();
		});
	}
	componentWillMount() {
		let url = "https://2020.elan.org.in/level/" + localStorage.getItem("id_token");
		fetch(url).then(response => response.json())
			.then(result => {
				this.setState({ url: result.URL });
				this.setState({ level: result.Level });
			});
	}
	render() {
		console.log(this.state.url);
		return (
			<div className="level-form won">
				<br />
				<img src={this.state.url} class="level-image" />
				<form onSubmit={this.handleSubmit}>
					<input type="name" className="answerTextbox" value={this.state.value} onChange={this.handleChange} />
					<br /><br />
					<input type="submit" className="answer-button" value="Submit" />
				</form>
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
	}
}
class LevelText extends React.Component {
	constructor(props) {
		super(props);
		this.state = { value: "", url: "", level: -3 };
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}
	handleSubmit(event) {
		event.preventDefault();
		let url = "https://2020.elan.org.in/answer/" + localStorage.getItem("id_token") + "/" + this.state.level.toString() + "/" + this.state.value;
		fetch(url).then(() => {
			window.location.reload();
		});
	}
	componentWillMount() {
		let url = "https://2020.elan.org.in/level/" + localStorage.getItem("id_token");
		fetch(url).then(response => response.json())
			.then(result => {
				this.setState({ url: result.URL });
				this.setState({ level: result.Level });
			});
	}
	render() {
		let strings = [];
		for (var i = 0; i < 13; i++) {
			strings[i] = this.state.url.substring(54 * i, 54 * i + 53);
		}
		return (
			<div className="level-form won">
				<br />
				<p className="white-text">{strings[0]}</p>
				<p className="white-text">{strings[1]}</p>
				<p className="white-text">{strings[2]}</p>
				<p className="white-text">{strings[3]}</p>
				<p className="white-text">{strings[4]}</p>
				<p className="white-text">{strings[5]}</p>
				<p className="white-text">{strings[6]}</p>
				<p className="white-text">{strings[7]}</p>
				<p className="white-text">{strings[8]}</p>
				<p className="white-text">{strings[9]}</p>
				<p className="white-text">{strings[10]}</p>
				<p className="white-text">{strings[11]}</p>
				<p className="white-text">{strings[12]}</p>
				<form onSubmit={this.handleSubmit}>
					<input type="name" className="answerTextbox" value={this.state.value} onChange={this.handleChange} />
					<br /><br />
					<input type="submit" className="answer-button" value="Submit" />
				</form>
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
	}
}

class LevelMidi extends React.Component {
	constructor(props) {
		super(props);
		this.state = { value: "", url: "", level: -3 };
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}
	handleSubmit(event) {
		event.preventDefault();
		let url = "https://2020.elan.org.in/answer/" + localStorage.getItem("id_token") + "/" + this.state.level.toString() + "/" + this.state.value;
		fetch(url).then(() => {
			window.location.reload();
		});
	}
	componentWillMount() {
		let url = "https://2020.elan.org.in/level/" + localStorage.getItem("id_token");
		fetch(url).then(response => response.json())
			.then(result => {
				this.setState({ url: result.URL });
				this.setState({ level: result.Level });
			});
	}
	render() {
		return (
			<div className="level-form won">
				<br />
				<a href={this.state.url}>Download</a>
				<form onSubmit={this.handleSubmit}>
					<input type="name" className="answerTextbox" value={this.state.value} onChange={this.handleChange} />
					<br /><br />
					<input type="submit" className="answer-button" value="Submit" />
				</form>
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
	}
}

class LevelUsername extends React.Component {
	constructor(props) {
		super(props);
		this.state = { value: "" };
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}
	handleSubmit(event) {
		event.preventDefault();
		let url = "https://2020.elan.org.in/doesUsernameExist/" + this.state.value ;
		fetch(url).then(response => response.json())
			.then(result => {
				if (result.message == "true") {
					alert("That username exists");
				}
				else {
					var loginUrl = "/adduser/" + JSON.parse(localStorage.getItem("email")).email + "/" + this.state.value + "/" + localStorage.getItem("id_token");
					fetch(loginUrl).then(() => {
						window.location.reload();
					});
				}
			});
	}
	render() {
		return (
			<div className="username-form won">
				<p>You are logged in, {JSON.parse(localStorage.getItem("email")).email}. </p>
				<p>Give us a username.</p>
				<form onSubmit={this.handleSubmit}>
					<input type="name" className="username" value={this.state.value} onChange={this.handleChange} />
					<br /><br />
					<input type="submit" className="username-button" value="Submit" />
				</form>
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
	}
}



class LevelRules extends React.Component {
	constructor(props) {
		super(props);
		this.handleAccepted = this.handleAccepted.bind(this);
	}
	logout() {
		localStorage.removeItem('id_token');
		localStorage.removeItem('access_token');
		localStorage.removeItem('profile');
		window.location.reload();
	}
	handleAccepted() {
		let url = "https://2020.elan.org.in/acceptedrules/" + (localStorage.getItem("id_token"));
		console.log(url);
		fetch(url);
		window.location.reload();
	}
	render() {
		return (
			<div className="rules-container won">
				<h1 className="rules">Rules</h1>
				<div class="rules" style={{textAlign: "left"}}>
					<div class="rules-content">
						<ol>
							<li>
								Mini Cryptex consists of 6 levels of increasing difficulty. You will receive successive questions upon solving and entering the answer of each level and/or completing the task involved.
                    </li>
							<li>
								The winner will be the first person to complete all levels. In case that no one is able to complete Cryptex before the 1800 hours 27th July, the person who occupies the first position on the leaderboard will be declared winner.
                    </li>
							<li>
								The competition is open only to freshers.
                    </li>
							<li>
								The questions, answers or any discussion related to them must not be posted anywhere. Doing so will result in disqualification. We will provide a platform to communicate with the moderators and ask for hints.
                    </li>
							<li>
								If a question involves typing answers in an answer box, the answers should be typed completely in lowercase, without any spaces, punctuation, special characters or numerals. The only exception to this is when the answer is the numeral greater than 1000 or when the answer is a special character.
                    </li>
							<li>
								If a question involves typing answers in an answer box, the answers should be typed completely in lowercase, without any spaces, punctuation, special characters or numerals. The only exception to this is when the answer is the numeral greater than 1000 or when the answer is a special character.
                    </li>
							<li>
								The questions may also involve interacting with the page itself, following hyperlinks etc.
                    </li>
							<li>
								Any attempt to access levels you have not yet reached or any form of attack on the servers will result in disqualification.
                    </li>
							<li>
								The organizer’s decision is final in any case.
                    </li>
						</ol>
					</div>
				</div>
				<button className="accept-rules" onClick={this.handleAccepted}>I accept</button>
				<br />
				<button onClick={this.logout}>Logout</button>
			</div>
		);
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
				<br />
				<div class="jumbotron animated fadeIn">
					<img src={image} class="main-image" />
					<p class="jumbotron-heading animated fadeIn">CRYPTEX</p>
					<p class="jumbotron-subtitle">Dive in and try out our mini online Treasure Hunt</p>
					<p class="jumbotron-subtitle">Online till 1800 hours, 27th July. </p>
					<button className="DiveInButton" onClick={this.authenticate}><div className="transform">D I V E &nbsp; I N</div></button>
				</div>
			</div>
		);
	}
}

class Callback extends React.Component {
	render() {
		return (<h1>Loading</h1>);
	}
}
class User {
	constructor(level, clientSecret, username) {
		this.level = level;
		this.clientSecret = clientSecret;
		this.username = username;
	}
	getLevel() {
		return (this.level);
	}
	getClientSecret() {
		return (this.clientSecret);
	}
	getUsername() {
		return (this.username);
	}
}


render(<App />, document.getElementById('app'));


if (module.hot) {
	module.hot.accept();
}