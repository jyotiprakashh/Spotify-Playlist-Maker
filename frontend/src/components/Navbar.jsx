import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate  } from 'react-router-dom';

function Navbar() {
  const [user, setUser] = useState(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const navigate = useNavigate(); 

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

  const handleLogout = async () => {
    try {
      const response = await axios.post('http://localhost:8080/logout', {}, { withCredentials: true });
      if (response.status === 200) {
        navigate('/'); // Use the navigate function to redirect to the root page
        console.log('User logged out');
      }
    } catch (error) {
      console.error('Error logging out:', error);
    }
  };


  return (
    <nav className="bg-spotifyblack p-4 flex justify-between items-center">
      <div className="text-spotifygreen text-xl font-bold ">SpotiMood</div>
      {user && (
        <div className="flex items-center gap-4">
          <div className="text-spotifywhite text-xl">{user.display_name}</div>
          <button onClick={() => setDropdownOpen(!dropdownOpen)} className="flex items-center">
            <img src={user.images || 'default-profile.png'} alt="Profile" className="w-8 h-8 rounded-full" />
          </button>
          {dropdownOpen && (
            <div className="absolute right-0 mt-20 w-48 bg-white rounded-md shadow-lg py-1 z-20">
              <button onClick={handleLogout} className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
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
