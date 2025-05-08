import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { GamesProvider } from './contexts/GamesContext';

// Pages
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import GamesPage from './pages/GamesPage';
import LobbyPage from './pages/LobbyPage';
import GamePlayPage from './pages/GamePlayPage';
import ProfilePage from './pages/ProfilePage';

// Protected route component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();
  
  if (isLoading) {
    return (
      <div className="min-h-screen bg-background-dark flex items-center justify-center">
        <div className="animate-pulse text-xl text-primary-500">Loading...</div>
      </div>
    );
  }
  
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

function App() {
  return (
    <Router>
      <AuthProvider>
        <GamesProvider>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            <Route path="/games" element={
              <ProtectedRoute>
                <GamesPage />
              </ProtectedRoute>
            } />
            
            <Route path="/lobby" element={
              <ProtectedRoute>
                <LobbyPage />
              </ProtectedRoute>
            } />
            
            <Route path="/play/:gameId" element={
              <ProtectedRoute>
                <GamePlayPage />
              </ProtectedRoute>
            } />
            
            <Route path="/profile" element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            } />
            
            {/* Fallback route */}
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </GamesProvider>
      </AuthProvider>
    </Router>
  );
}

export default App;