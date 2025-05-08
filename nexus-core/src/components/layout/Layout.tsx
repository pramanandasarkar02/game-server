import React from 'react';
import Header from './Header';
import BackgroundParticles from './BackgroundParticles';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col bg-background-dark text-white">
      <BackgroundParticles />
      <Header />
      <main className="flex-grow container mx-auto px-4 py-6">
        {children}
      </main>
      <footer className="bg-background-dark border-t border-gray-800 py-4">
        <div className="container mx-auto px-4 text-center text-gray-400 text-sm">
          <p>&copy; {new Date().getFullYear()} Nexus Core. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default Layout;