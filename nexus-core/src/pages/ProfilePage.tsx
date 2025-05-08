import React from 'react';
import { User, Trophy, Award, Flag, Clock, Zap } from 'lucide-react';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import { useAuth } from '../contexts/AuthContext';

const ProfilePage: React.FC = () => {
  const { user } = useAuth();

  if (!user) return null;

  // Calculate win rate
  const winRate = user.gamesPlayed > 0 
    ? Math.round((user.gamesWon / user.gamesPlayed) * 100) 
    : 0;

  // Calculate progress to next level (simplified)
  const xpForCurrentLevel = user.level * 1000;
  const xpForNextLevel = (user.level + 1) * 1000;
  const xpProgress = Math.min(100, Math.max(0, ((user.score - xpForCurrentLevel) / (xpForNextLevel - xpForCurrentLevel)) * 100));

  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Player Profile</h1>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="md:col-span-1">
            <Card className="p-6">
              <div className="flex flex-col items-center text-center">
                <div className="w-24 h-24 rounded-full bg-primary-700 flex items-center justify-center mb-4">
                  <User className="h-12 w-12 text-primary-100" />
                </div>
                <h2 className="text-xl font-bold">{user.username}</h2>
                <p className="text-gray-400 mb-4">{user.email}</p>
                
                <div className="w-full mt-2">
                  <div className="flex justify-between text-sm mb-1">
                    <span>Level {user.level}</span>
                    <span>Level {user.level + 1}</span>
                  </div>
                  <div className="w-full bg-gray-700 rounded-full h-2.5">
                    <div 
                      className="bg-primary-500 h-2.5 rounded-full" 
                      style={{ width: `${xpProgress}%` }}
                    ></div>
                  </div>
                  <p className="text-xs text-gray-400 mt-1">
                    {Math.round(xpProgress)}% to next level
                  </p>
                </div>
              </div>
              
              <div className="mt-6 grid grid-cols-2 gap-4">
                <div className="flex flex-col items-center p-3 bg-background-dark rounded-lg">
                  <Trophy className="h-5 w-5 text-yellow-500 mb-1" />
                  <span className="text-lg font-bold">{user.score}</span>
                  <span className="text-xs text-gray-400">Total Score</span>
                </div>
                <div className="flex flex-col items-center p-3 bg-background-dark rounded-lg">
                  <Award className="h-5 w-5 text-purple-500 mb-1" />
                  <span className="text-lg font-bold">{user.gamesWon}</span>
                  <span className="text-xs text-gray-400">Games Won</span>
                </div>
              </div>
            </Card>
          </div>

          <div className="md:col-span-2">
            <Card className="p-6 mb-6">
              <h2 className="text-xl font-bold mb-4">Stats Overview</h2>
              
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="flex items-center p-4 bg-background-dark rounded-lg">
                  <div className="mr-4 bg-green-900/30 p-3 rounded-full">
                    <Flag className="h-6 w-6 text-green-500" />
                  </div>
                  <div>
                    <div className="text-2xl font-bold">{user.gamesPlayed}</div>
                    <div className="text-sm text-gray-400">Games Played</div>
                  </div>
                </div>
                
                <div className="flex items-center p-4 bg-background-dark rounded-lg">
                  <div className="mr-4 bg-blue-900/30 p-3 rounded-full">
                    <Trophy className="h-6 w-6 text-blue-500" />
                  </div>
                  <div>
                    <div className="text-2xl font-bold">{winRate}%</div>
                    <div className="text-sm text-gray-400">Win Rate</div>
                  </div>
                </div>
                
                <div className="flex items-center p-4 bg-background-dark rounded-lg">
                  <div className="mr-4 bg-purple-900/30 p-3 rounded-full">
                    <Zap className="h-6 w-6 text-purple-500" />
                  </div>
                  <div>
                    <div className="text-2xl font-bold">{user.level}</div>
                    <div className="text-sm text-gray-400">Current Level</div>
                  </div>
                </div>
                
                <div className="flex items-center p-4 bg-background-dark rounded-lg">
                  <div className="mr-4 bg-orange-900/30 p-3 rounded-full">
                    <Clock className="h-6 w-6 text-orange-500" />
                  </div>
                  <div>
                    <div className="text-2xl font-bold">5h 23m</div>
                    <div className="text-sm text-gray-400">Time Played</div>
                  </div>
                </div>
              </div>
            </Card>
            
            <Card className="p-6">
              <h2 className="text-xl font-bold mb-4">Recent Achievements</h2>
              
              <div className="space-y-4">
                <div className="flex items-center justify-between p-3 bg-background-dark rounded-lg border-l-4 border-yellow-500">
                  <div className="flex items-center">
                    <Award className="h-5 w-5 text-yellow-500 mr-3" />
                    <div>
                      <div className="font-medium">First Victory</div>
                      <div className="text-xs text-gray-400">Won your first game</div>
                    </div>
                  </div>
                  <div className="text-xs text-gray-400">3 days ago</div>
                </div>
                
                <div className="flex items-center justify-between p-3 bg-background-dark rounded-lg border-l-4 border-blue-500">
                  <div className="flex items-center">
                    <Trophy className="h-5 w-5 text-blue-500 mr-3" />
                    <div>
                      <div className="font-medium">Puzzle Master</div>
                      <div className="text-xs text-gray-400">Completed puzzle in under 30 moves</div>
                    </div>
                  </div>
                  <div className="text-xs text-gray-400">5 days ago</div>
                </div>
                
                <div className="flex items-center justify-between p-3 bg-background-dark rounded-lg border-l-4 border-green-500">
                  <div className="flex items-center">
                    <Zap className="h-5 w-5 text-green-500 mr-3" />
                    <div>
                      <div className="font-medium">Level Up</div>
                      <div className="text-xs text-gray-400">Reached level 5</div>
                    </div>
                  </div>
                  <div className="text-xs text-gray-400">1 week ago</div>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default ProfilePage;