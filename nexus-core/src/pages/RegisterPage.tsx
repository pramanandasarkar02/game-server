import React from 'react';
import { Link } from 'react-router-dom';
import Layout from '../components/layout/Layout';
import RegisterForm from '../components/auth/RegisterForm';

const RegisterPage: React.FC = () => {
  return (
    <Layout>
      <div className="min-h-[calc(100vh-200px)] flex flex-col items-center justify-center">
        <div className="w-full max-w-md">
          <RegisterForm />
          <div className="text-center mt-4 text-gray-400">
            Already have an account?{' '}
            <Link to="/login" className="text-primary-400 hover:text-primary-300">
              Login instead
            </Link>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default RegisterPage;