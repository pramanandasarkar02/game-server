import React from 'react';
import { Link } from 'react-router-dom';
import Layout from '../components/layout/Layout';
import LoginForm from '../components/auth/LoginForm';

const LoginPage: React.FC = () => {
  return (
    <Layout>
      <div className="min-h-[calc(100vh-200px)] flex flex-col items-center justify-center">
        <div className="w-full max-w-md">
          <LoginForm />
          <div className="text-center mt-4 text-gray-400">
            Don't have an account?{' '}
            <Link to="/register" className="text-primary-400 hover:text-primary-300">
              Create one now
            </Link>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default LoginPage;