import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { TowerControl as GameController, User, LogOut } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import Button from '../ui/Button';

const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="bg-background-dark border-b border-gray-800">
      <div className="container mx-auto px-4 py-3 flex justify-between items-center">
        <Link to="/" className="flex items-center space-x-2">
          <GameController className="h-8 w-8 text-primary-500" />
          <h1 className="text-xl md:text-2xl font-bold text-white">
            <span className="text-primary-500">Nexus</span> Core
          </h1>
        </Link>

        <nav>
          <ul className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <li className="hidden md:block">
                  <Link to="/games" className="text-gray-300 hover:text-white transition-colors">
                    Games
                  </Link>
                </li>
                <li className="hidden md:block">
                  <Link to="/profile" className="text-gray-300 hover:text-white transition-colors">
                    Profile
                  </Link>
                </li>
                <li className="flex items-center space-x-2">
                  <Link to="/profile" className="flex items-center space-x-2 px-3 py-1 rounded-full bg-background-light">
                    <User className="h-5 w-5 text-primary-400" />
                    <span className="text-sm font-medium">{user?.username}</span>
                    <span className="text-xs px-2 py-1 rounded-full bg-primary-700 text-primary-100">Lv.{user?.level}</span>
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="p-2 text-gray-400 hover:text-white transition-colors" 
                    aria-label="Logout"
                  >
                    <LogOut className="h-5 w-5" />
                  </button>
                </li>
              </>
            ) : (
              <>
                <li>
                  <Link to="/login">
                    <Button variant="outline" size="sm">
                      Login
                    </Button>
                  </Link>
                </li>
                <li>
                  <Link to="/register">
                    <Button size="sm">
                      Register
                    </Button>
                  </Link>
                </li>
              </>
            )}
          </ul>
        </nav>
      </div>
    </header>
  );
};

export default Header;