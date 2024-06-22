import React, { useState, useEffect } from 'react';
import axios from 'axios';

function MoodInput() {
  const [mood, setMood] = useState('');
  const [playlist, setPlaylist] = useState([]);
  const [playlistURL, setPlaylistURL] = useState('');

  const moods = [
    "happy", "sad", "relaxed", "energetic", "angry", 
    "romantic", "party", "focused", "nostalgic", "melancholic"
  ];

  const [user, setUser] = useState(null);

  useEffect(() => {
    async function fetchUserProfile() {
      try {
        const response = await axios.get('http://localhost:8080/api/user-profile', { withCredentials: true });
        setUser(response.data);
        // console.log('User profile:', response.data);
      } catch (error) {
        console.error('Error fetching user profile:', error);
      }
    }

    fetchUserProfile();
  }, []);



  const handleGeneratePlaylist = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post('http://localhost:8080/api/generate-playlist', { mood });
      console.log('Playlist response:', response.data.playlist);
      setPlaylist(response.data.playlist);
      setPlaylistURL(''); // Clear the playlist URL when generating a new playlist
    } catch (error) {
      console.error('Error generating playlist:', error);
    }
  };

  const handleAddPlaylist = async () => {
    const playlistName = `${mood.charAt(0).toUpperCase() + mood.slice(1)} Playlist`;
    try {
      const response = await axios.post('http://localhost:8080/api/addPlaylist', {playlist: playlist.map(track => track.name), name: playlistName});
      console.log('Playlist URL response:', response.data.playlist_url);
      setPlaylistURL(response.data.playlist_url);
    } catch (error) {
      console.error('Error adding playlist:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-r bg-spotifyblack text-white flex flex-col items-center justify-center p-4">
      <div>
      <div className='text-center'>
        {user ? (
          <h1 className="text-3xl font-bold mb-4">Hi, {user.display_name}!</h1>
        ) : (
          <h1 className="text-3xl font-bold mb-4">Loading...</h1>
        )}
      </div>
        <h1 className="text-2xl font-medium mb-4">How are you feeling today?</h1>
      </div>
    <form onSubmit={handleGeneratePlaylist} className="w-full max-w-md bg-white p-8 rounded-lg shadow-md text-black">
    <input
          type="text"
          value={mood}
          onChange={(e) => setMood(e.target.value)}
          placeholder="Enter your mood"
          className="w-full p-2 mb-4 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600"
        />
     
      <h1 className="text-2xl font-bold text-center mb-4">Generate Your Mood Playlist</h1>
      {/* <select value={mood} onChange={(e) => setMood(e.target.value)} className="w-full p-2 mb-4 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600">
        <option value="">Select your mood</option>
        {moods.map((m, index) => (
          <option key={index} value={m}>{m}</option>
        ))}
      </select> */}
      <button type="submit" className="w-full bg-spotifygreen text-white py-2 rounded-lg font-semibold hover:bg-green-700 transition duration-300">Generate Playlist</button>
    </form>

    {playlist.length > 0 && (
      <div className="w-full max-w-2xl bg-white p-8 mt-8 rounded-lg shadow-md text-black">
        <h2 className="text-xl font-bold mb-4">Generated Playlist</h2>
        <ul className="space-y-4">
          {playlist.map((track, index) => (
            <li key={index} className="flex items-center space-x-4">
              <img src={track.image} alt={`${track.name} cover`} className="w-12 h-12 rounded-lg shadow-md" />
              <span className="font-medium">{track.name}</span>
            </li>
          ))}
        </ul>
        <button onClick={handleAddPlaylist} className="w-full bg-green-600 text-white py-2 mt-4 rounded-lg font-semibold hover:bg-green-700 transition duration-300">Add Playlist to Spotify</button>
      </div>
    )}

    {playlistURL && (
      <div className="w-full max-w-2xl bg-white p-8 mt-8 rounded-lg shadow-md text-black">
        <h2 className="text-xl font-bold mb-4">Spotify Playlist</h2>
        <a href={playlistURL} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Open Playlist in Spotify</a>
      </div>
    )}
  </div>
  );
}

export default MoodInput;
