import React from "react";

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

class NavbarLoggedIn extends React.Component {
    logout() {
        localStorage.removeItem('id_token');
        localStorage.removeItem('access_token');
        localStorage.removeItem('profile');
        window.location.reload();
    }
    displayLevels() {
        console.log("Hi");
    }
    render() {
        return (
            <nav class="animated fadeInDown">
                <ul>
                    <a href="#" onClick={this.displayLevels}><li>Levels</li></a>
                    <a href="/rules"><li>Guidelines</li></a>
                    <a href="/" className="main animated flipInX"><li>C R Y P T E X</li></a>
                    <a href="/leaderboardtable" id="leaderboard-nav-link"><li>Leaderboard</li></a>
                    <a href="#" onClick={this.logout}><li>Logout</li></a>
                </ul>
            </nav>
        );
    }
}

export {
    Navbar,
    NavbarLoggedIn
}