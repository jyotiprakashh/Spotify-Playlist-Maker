import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function Navbar() {
  const [user, setUser] = useState(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
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

  const handleLogout = async () => {
    try {
      const response = await axios.post(
        "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/logout",
        {},
        { withCredentials: true }
      );
      if (response.status === 200) {
        setUser(null);
        sessionStorage.clear();
        console.log("User logged out");
        window.location.href = "https://accounts.spotify.com/logout";
      }
    } catch (error) {
      console.error("Error logging out:", error);
    }
  };

  return (
    <nav className=" p-4  top-0 left-0 w-full bg-opacity-100  ">
      <div className="max-w-7xl mx-auto flex justify-between items-center">
        <div className="text-spotifygreen text-2xl font-bold">SpotiMood</div>

        <div className="hidden md:flex gap-8 items-center">
          <a href="/" className="text-white text-md hover:text-spotifygreen">
            Home
          </a>
          {user && (
            <div className="relative flex items-center gap-4">
              <span className="text-white">{user.display_name}</span>
              <button
                onClick={() => setDropdownOpen(!dropdownOpen)}
                className="flex items-center gap-2"
              >
                <img
                  src={user.images || "default-profile.png"}
                  alt="Profile"
                  className="w-8 h-8 rounded-full"
                />
              </button>
              {dropdownOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-2 z-20">
                  <button
                    onClick={handleLogout}
                    className="block w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100"
                  >
                    Logout
                  </button>
                </div>
              )}
            </div>
          )}
        </div>

        <div className="md:hidden flex items-center gap-4">
          {user && (
            <div className="flex items-center gap-2">
              <img
                src={user.images || "default-profile.png"}
                alt="Profile"
                className="w-8 h-8 rounded-full"
              />
              <span className="text-white">{user.display_name}</span>
            </div>
          )}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="text-white focus:outline-none"
          >
            {mobileMenuOpen ? (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M6 18 18 6M6 6l12 12"
                />
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                />
              </svg>
            )}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-spotifyblack text-white py-4 space-y-4">
          <a
            href="/"
            className="block text-center text-md hover:text-spotifygreen"
          >
            Home
          </a>
          {user && (
            <div className="flex flex-col items-center">
              <button
                onClick={handleLogout}
                className="mt-4 px-4 py-2 text-gray-700 bg-white rounded-lg hover:bg-gray-100"
              >
                Logout
              </button>
            </div>
          )}
        </div>
      )}
    </nav>
  );
}

export default Navbar;
