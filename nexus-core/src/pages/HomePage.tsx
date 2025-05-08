import React from 'react';
import { Link } from 'react-router-dom';
import { TowerControl as GameController, Users, Trophy, Zap } from 'lucide-react';
import Layout from '../components/layout/Layout';
import Button from '../components/ui/Button';
import { useAuth } from '../contexts/AuthContext';

const HomePage: React.FC = () => {
  const { isAuthenticated } = useAuth();

  return (
    <Layout>
      <div className="min-h-[calc(100vh-200px)] flex flex-col items-center justify-center text-center">
        <div className="max-w-4xl">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Welcome to <span className="text-primary-500">Nexus</span> Core
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 mb-10">
            The ultimate gaming platform where players connect, compete, and conquer.
          </p>

          <div className="flex flex-wrap justify-center gap-4 mb-12">
            {isAuthenticated ? (
              <Link to="/games">
                <Button size="lg">
                  <GameController className="mr-2 h-5 w-5" />
                  Play Now
                </Button>
              </Link>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="outline" size="lg">
                    Login
                  </Button>
                </Link>
                <Link to="/register">
                  <Button size="lg">
                    Create Account
                  </Button>
                </Link>
              </>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-10">
            <div className="bg-background-card p-6 rounded-lg shadow-lg">
              <div className="bg-primary-700/30 w-14 h-14 rounded-full flex items-center justify-center mb-4 mx-auto">
                <GameController className="h-8 w-8 text-primary-400" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Multiple Games</h3>
              <p className="text-gray-400">
                Enjoy a variety of challenging games from classic Tic Tac Toe to mind-bending puzzles.
              </p>
            </div>

            <div className="bg-background-card p-6 rounded-lg shadow-lg">
              <div className="bg-secondary-700/30 w-14 h-14 rounded-full flex items-center justify-center mb-4 mx-auto">
                <Users className="h-8 w-8 text-secondary-400" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Play with Others</h3>
              <p className="text-gray-400">
                Join game lobbies and play with other players in real-time matches.
              </p>
            </div>

            <div className="bg-background-card p-6 rounded-lg shadow-lg">
              <div className="bg-accent-700/30 w-14 h-14 rounded-full flex items-center justify-center mb-4 mx-auto">
                <Trophy className="h-8 w-8 text-accent-400" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Level Up</h3>
              <p className="text-gray-400">
                Earn points, increase your level, and climb the leaderboards with each victory.
              </p>
            </div>
          </div>
        </div>

        <div className="mt-16 w-full max-w-4xl bg-gradient-to-r from-primary-900/50 to-secondary-900/50 rounded-lg p-6 flex flex-col md:flex-row items-center justify-between">
          <div>
            <h3 className="text-xl font-semibold mb-2 flex items-center">
              <Zap className="h-5 w-5 mr-2 text-yellow-400" />
              Ready to Play?
            </h3>
            <p className="text-gray-300">Join now and start your gaming journey!</p>
          </div>
          <div className="mt-4 md:mt-0">
            {isAuthenticated ? (
              <Link to="/games">
                <Button variant="accent">
                  Browse Games
                </Button>
              </Link>
            ) : (
              <Link to="/register">
                <Button variant="accent">
                  Create Free Account
                </Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default HomePage;