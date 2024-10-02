import React from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const LandingPage = () => {
  const navigate = useNavigate();

  const handleGetStarted = async () => {
    try {
      const response = await axios.get("https://racial-angela-jyotiprakashh-5b571460.koyeb.app/api/getUserID", {
        withCredentials: true,
      });
      if (response.data.user_id) {
        navigate("/app");
      }
    } catch (error) {
      window.location.href = "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/login";
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-900 via-spotifyblack to-spotifyblack text-white overflow-hidden">
      <div className="relative container mx-auto px-4 py-16">
        <div className="absolute top-0 left-0 w-72 h-72 bg-spotifygreen rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
        <div className="absolute top-0 right-0 w-72 h-72 bg-spotifygreen rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
        
        <nav className="relative flex justify-between items-center mb-16">
          <div className="text-3xl font-extrabold text-white">SpotiMood</div>
        </nav>

        <div className="relative flex flex-col items-center justify-center mt-20">
          <h1 className="text-5xl md:text-6xl lg:text-7xl font-extrabold text-center mb-6 bg-clip-text text-transparent bg-gradient-to-r from-spotifygreen to-emerald-400">
            Discover Your Playlist Mood
          </h1>
          <p className="text-lg md:text-xl lg:text-2xl text-center mb-12 max-w-3xl px-4">
            Let SpotiMood generate the perfect Spotify playlist based on your
            emotions. AI-driven playlist curation for every moment.
          </p>

          <div className="flex flex-col sm:flex-row gap-4">
            <button
              onClick={handleGetStarted}
              className="px-8 py-3 bg-gradient-to-r from-spotifygreen/80 to-spotifygreen text-white font-semibold rounded-full hover:from-spotifygreen hover:to-spotifygreen transform hover:scale-105 transition duration-300 ease-in-out shadow-lg"
            >
              Get Started
            </button>
            <a
              href="#learn-more"
              className="px-8 py-3 bg-white bg-opacity-20 text-white font-semibold rounded-full hover:bg-opacity-30 transform hover:scale-105 transition duration-300 ease-in-out shadow-lg"
            >
              Learn More
            </a>
          </div>
        </div>

        <footer className="relative mt-24 text-center text-gray-400 bottom-0 ">
          <p>&copy; 2024 SpotiMood. By Jyoti Prakash</p>
        </footer>
      </div>
    </div>
  );
};

export default LandingPage;
