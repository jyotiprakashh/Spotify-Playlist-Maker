# 🎵 Spotify Mood-Based Playlist Generator

## Overview
This project is a web application that interacts with the Spotify API to generate mood-based playlists. Users can log in with their Spotify accounts, input their current mood, and receive a playlist tailored to that mood. The application is built using **Go** for the backend and **React Vite** for the frontend.

## Features
- User Authentication: Log in with your Spotify account.
- Mood-Based Playlists: Generate playlists based on your current mood.
- User Profile: Fetch and display user profile information from Spotify.
- Playlist Creation: Create and add tracks to a new Spotify playlist.

# Technology Stack
- Backend: Go
- Frontend: React Vite
- API: Spotify Web API
- OAuth: Spotify OAuth 2.0
- Containerization: Docker

# Prerequisites
- Go 1.20+
- Node.js 14+
- Spotify Developer Account (Client ID and Secret)
- Docker
  
# Installation
**Backend**
1. Clone the repository:
```bash
git clone https://github.com/jyotiprakashh/Spotify-Playlist-Maker
cd spotify-mood-playlist
```

2. Go to the backend folder
```bash
cd backend
```
3. Set up environment variables:
```bash
SPOTIFY_ID=your_spotify_client_id
SPOTIFY_SECRET=your_spotify_client_secret
PORT=8080
```

4. Change the Redirect URL in the main file to your redirect URL.

5. Run the application using Docker:
```bash
docker build -t spotify-mood-backend .
docker run -d --restart=always  -p 8080:8080 playlist-maker
```

**Frontend**
1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the frontend application:
```bash
npm run dev
```


## API Endpoints
- GET /login: Redirects to Spotify login for user authentication.
- GET /callback: Callback URL for Spotify OAuth.
- POST /api/generate-playlist: Generates a playlist based on the provided mood.
- POST /api/addPlaylist: Creates a new playlist and adds tracks to it.
- GET /api/user-profile: Fetches the user's Spotify profile information.
- POST /logout: Logs out the user.
