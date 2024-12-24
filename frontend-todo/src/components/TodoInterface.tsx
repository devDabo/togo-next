import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CardComponent from './CardComponent';

interface Todo {
  id: number;
  todo: string;
}

interface TodoInterfaceProps {
  backendName: string;
}

const TodoInterface: React.FC<TodoInterfaceProps> = ({ backendName }) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState({ todo: '' });
  const [updateTodo, setUpdateTodo] = useState({ id: '', todo: '' });

  const backgroundColors: { [key: string]: string } = {
    go: 'bg-cyan-500',
  };

  const buttonColors: { [key: string]: string } = {
    go: 'bg-cyan-700 hover:bg-blue-600',
  };

  const bgColor = backgroundColors[backendName as keyof typeof backgroundColors] || 'bg-gray-200';
  const btnColor = buttonColors[backendName as keyof typeof buttonColors] || 'bg-gray-500 hover:bg-gray-600';

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/${backendName}/todo`);
        setTodos(response.data.reverse());
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [backendName, apiUrl]);

  const createTodo = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${apiUrl}/api/${backendName}/todo`, newTodo);
      setTodos([response.data, ...todos]);
      setNewTodo({ todo: '' });
    } catch (error) {
      console.error('Error creating todo:', error);
    }
  };

  const handleUpdateTodo = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      await axios.put(`${apiUrl}/api/${backendName}/todo/${updateTodo.id}`, { todo: updateTodo.todo });
      setUpdateTodo({ id: '', todo: '' });
      setTodos(
        todos.map((todo) => {
          if (todo.id === parseInt(updateTodo.id)) {
            return { ...todo, todo: updateTodo.todo };
          }
          return todo;
        })
      );
    } catch (error) {
      console.error('Error updating todo:', error);
    }
  };

  const deleteTodo = async (todoId: number) => {
    try {
      await axios.delete(`${apiUrl}/api/${backendName}/todo/${todoId}`);
      setTodos(todos.filter((todo) => todo.id !== todoId));
    } catch (error) {
      console.error('Error deleting todo:', error);
    }
  };

  return (
    <div className={`todo-interface ${bgColor} ${backendName} w-full max-w-md p-4 my-4 rounded shadow`}>
      <h2 className="text-xl font-bold text-center text-white mb-6">{`Todo app`}</h2>

      <form onSubmit={createTodo} className="mb-6 p-4 bg-blue-100 rounded shadow">
        <input
          placeholder="Add a new todo"
          value={newTodo.todo}
          onChange={(e) => setNewTodo({ todo: e.target.value })}
          className="mb-2 w-full p-2 border border-gray-300 rounded"
        />
        <button type="submit" className="w-full p-2 text-white bg-blue-500 rounded hover:bg-blue-600">
          Add Todo
        </button>
      </form>

      <form onSubmit={handleUpdateTodo} className="mb-6 p-4 bg-blue-100 rounded shadow">
        <input
          placeholder="Todo Id"
          value={updateTodo.id}
          onChange={(e) => setUpdateTodo({ ...updateTodo, id: e.target.value })}
          className="mb-2 w-full p-2 border border-gray-300 rounded"
        />
        <input
          placeholder="New Todo"
          value={updateTodo.todo}
          onChange={(e) => setUpdateTodo({ ...updateTodo, todo: e.target.value })}
          className="mb-2 w-full p-2 border border-gray-300 rounded"
        />
        <button type="submit" className="w-full p-2 text-white bg-green-500 rounded hover:bg-green-600">
          Update Todo
        </button>
      </form>

      <div className="space-y-4">
        {todos.map((todo) => (
          <div key={todo.id} className="flex items-center justify-between bg-white p-4 rounded-lg shadow">
            <CardComponent todo={todo} />
            <button onClick={() => deleteTodo(todo.id)} className={`${btnColor} text-white py-2 px-4 rounded`}>
              Delete Todo
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TodoInterface;
