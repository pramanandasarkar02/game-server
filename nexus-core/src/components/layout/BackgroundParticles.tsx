import React, { useEffect } from 'react';

interface Particle {
  id: number;
  x: number;
  y: number;
  size: number;
  speed: number;
  opacity: number;
  color: string;
}

const BackgroundParticles: React.FC = () => {
  useEffect(() => {
    // Number of particles to create
    const particleCount = 20;
    const particles: Particle[] = [];
    const colors = ['#7c3aed', '#00e6af', '#f97316'];
    
    // Create container
    const container = document.querySelector('.particles-container');
    if (!container) return;
    
    // Clear existing particles
    container.innerHTML = '';
    
    // Create particles
    for (let i = 0; i < particleCount; i++) {
      const particle: Particle = {
        id: i,
        x: Math.random() * 100, // percentage
        y: Math.random() * 100, // percentage
        size: Math.random() * 15 + 5, // size in pixels
        speed: Math.random() * 3 + 1, // animation duration modifier
        opacity: Math.random() * 0.3 + 0.1, // opacity value
        color: colors[Math.floor(Math.random() * colors.length)]
      };
      
      particles.push(particle);
      
      // Create particle element
      const element = document.createElement('div');
      element.className = 'particle';
      element.style.left = `${particle.x}%`;
      element.style.top = `${particle.y}%`;
      element.style.width = `${particle.size}px`;
      element.style.height = `${particle.size}px`;
      element.style.backgroundColor = particle.color;
      element.style.opacity = particle.opacity.toString();
      element.style.animationDuration = `${particle.speed * 3}s`;
      element.style.animationDelay = `${Math.random() * 5}s`;
      
      container.appendChild(element);
    }
    
    // Clean up
    return () => {
      if (container) {
        container.innerHTML = '';
      }
    };
  }, []);
  
  return <div className="particles-container" />;
};

export default BackgroundParticles;