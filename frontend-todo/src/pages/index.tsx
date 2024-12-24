
import React from 'react';
import TodoInterface from '@/components/TodoInterface';

const Home: React.FC = () => {
  return (
    <main className="flex flex-wrap justify-center items-start min-h-screen bg-gray-100">
      <div className="m-4">
        <TodoInterface backendName="go" />
      </div>
    </main>
  );
};

export default Home;