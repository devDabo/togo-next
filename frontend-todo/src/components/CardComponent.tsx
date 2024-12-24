import React from 'react';

interface Todo {
  id: number;
  todo: string;
}

const CardComponent: React.FC<{ todo: Todo }> = ({ todo }) => {
  return (
    <div className="bg-white shadow-lg rounded-lg p-2 mb-2 hover:bg-gray-100">
      <div className="text-sm text-gray-600">Id: {todo.id}</div>
      <div className="text-lg font-semibold text-gray-800">{todo.todo}</div>
    </div>
  );
};

export default CardComponent;
