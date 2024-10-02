import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Navbar from './Navbar';


function MoodInput() {
  const [mood, setMood] = useState("");
  const [playlist, setPlaylist] = useState([]);
  const [playlistURL, setPlaylistURL] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedTrack, setSelectedTrack] = useState(null);
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    async function getUserid() {
      try {
        const response = await axios.get(
          "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/api/getUserID",
          { withCredentials: true }
        );
        if (!response.data.user_id) {
          console.log("User ID not found");
          navigate("/");
        }
        console.log("User ID:", response.data.user_id);
        localStorage.setItem("userID", response.data.user_id);
      } catch (error) {
        console.error("Error fetching user profile:", error);
      }
    }
    getUserid();
  }, [navigate]);

  useEffect(() => {
    async function fetchUserProfile() {
      try {
        const response = await axios.get(
          "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/api/user-profile",
          { withCredentials: true }
        );
        setUser(response.data);
      } catch (error) {
        console.error("Error fetching user profile:", error);
      }
    }
    fetchUserProfile();
  }, []);

  const handleGeneratePlaylist = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await axios.post(
        "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/api/generate-playlist",
        { mood }
      );
      setPlaylist(response.data.playlist);
      setPlaylistURL("");
      console.log("Playlist:", response.data.playlist);
    } catch (error) {
      console.error("Error generating playlist:", error);
    }
    setLoading(false);
  };

  const handleAddPlaylist = async () => {
    const playlistName = `${
      mood.charAt(0).toUpperCase() + mood.slice(1)
    } Playlist`;
    try {
      const response = await axios.post(
        "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/api/addPlaylist",
        { playlist: playlist.map((track) => track.name), name: playlistName }
      );
      setPlaylistURL(response.data.playlist_url);
    } catch (error) {
      console.error("Error adding playlist:", error);
    }
  };

  const handlePlayTrack = (track) => {
    setSelectedTrack(track);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-900 via-spotifyblack to-spotifyblack text-white p-2">
      <Navbar className=""  />

      <div className="max-w-6xl mx-auto mt-5">
        <header className="text-center mb-12">
          <h1 className="md:text-5xl text-3xl font-bold mb-4">
            Hey {user ? user.display_name : "there"}!
          </h1>
          <p className="md:text-xl text-md text-gray-300 ">
            Let's create a playlist that matches your mood!
          </p>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:p-0 p-3">
          <div className="bg-white bg-opacity-10 backdrop-filter backdrop-blur-lg rounded-xl shadow-lg p-6 md:h-[500px] h-[400px] flex flex-col">
            <h2 className="text-2xl font-semibold mb-4">How was your day?</h2>
            <form
              onSubmit={handleGeneratePlaylist}
              className="flex-grow flex flex-col"
            >
              <textarea
                value={mood}
                onChange={(e) => setMood(e.target.value)}
                placeholder="Tell us about your day in detail..."
                className="w-full p-3 mb-4 bg-white bg-opacity-20 border-none rounded-lg focus:ring-2 focus:ring-spotifygreen focus:outline-none transition duration-300 placeholder-gray-400 text-white resize-none flex-grow"
              />
              <button
                type="submit"
                className="w-full bg-green-500 hover:bg-green-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center"
                disabled={loading}
              >
                {loading ? (
                  <svg
                    className="animate-spin h-5 w-5 mr-3"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                ) : (
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M14 5l7 7m0 0l-7 7m7-7H3"
                    />
                  </svg>
                )}
                {loading ? "Generating..." : "Generate Playlist"}
              </button>
            </form>
          </div>

          <div className="bg-white bg-opacity-10 backdrop-filter backdrop-blur-lg rounded-xl shadow-lg p-6 h-[500px] flex flex-col mb-20">
            <h2 className="text-2xl font-semibold mb-4">Your Mood Playlist</h2>
            {playlist.length > 0 ? (
              <div className="flex-grow flex flex-col">
                <div className="overflow-y-auto flex-grow mb-4 h-[300px]">
                  <ul className="space-y-4">
                    {playlist.map((track, index) => (
                      <li
                        key={index}
                        className="flex items-center space-x-4 hover:bg-white hover:bg-opacity-20 p-2 rounded-lg cursor-pointer transition duration-300"
                        onClick={() => handlePlayTrack(track)}
                      >
                        <img
                          src={track.image}
                          alt={`${track.name} cover`}
                          className="w-12 h-12 rounded-lg shadow-md"
                        />
                        <span className="font-medium">{track.name}</span>
                      </li>
                    ))}
                  </ul>
                </div>
                <button
                  onClick={handleAddPlaylist}
                  className="w-full bg-green-500 hover:bg-green-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center"
                >
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"
                    />
                  </svg>
                  Add Playlist to Spotify
                </button>
              </div>
            ) : (
              <p className="text-center text-gray-400 flex-grow flex items-center justify-center">
                Your playlist will appear here after generation.
              </p>
            )}
            {playlistURL && (
              <a
                href={playlistURL}
                target="_blank"
                rel="noopener noreferrer"
                className="mt-4 inline-flex items-center text-green-400 hover:underline"
              >
                <svg
                  className="w-5 h-5 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                  />
                </svg>
                Open Playlist in Spotify
              </a>
            )}
          </div>
        </div>

        {selectedTrack && (
          <div className="fixed bottom-0 left-0 w-full bg-gray-900 bg-opacity-90 text-white p-4">
            <div className="flex items-center justify-between max-w-6xl mx-auto">
              <div className="flex items-center space-x-4">
                <img
                  src={selectedTrack.image}
                  alt={`${selectedTrack.name} cover`}
                  className="md:w-16 md:h-16 w-10 h-10 rounded-lg shadow-md "
                />
                <div>
                  <h3 className="md:text-xl text-md font-bold">{selectedTrack.name}</h3>
                  <a
                    href={selectedTrack.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-green-400 text-sm hover:underline"
                  >
                    Listen on Spotify
                  </a>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default MoodInput;
