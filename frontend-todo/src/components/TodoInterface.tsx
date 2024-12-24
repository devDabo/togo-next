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
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);

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
      setLoading(true);
      setError('');
      try {
        const response = await axios.get(`${apiUrl}/api/${backendName}/todo`);
        setTodos(response.data ? response.data.reverse() : []);
      } catch (error) {
        setError('Failed to fetch todos. Please try again later.');
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [backendName, apiUrl]);

  const createTodo = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    if (!newTodo.todo.trim()) {
      setError('Todo cannot be empty');
      return;
    }
    
    try {
      const response = await axios.post(`${apiUrl}/api/${backendName}/todo`, newTodo);
      setTodos([response.data, ...todos]);
      setNewTodo({ todo: '' });
    } catch (error: unknown) {
      setError(error.response?.data?.message || 'Failed to create todo');
      console.error('Error creating todo:', error);
    }
  };

  const handleUpdateTodo = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    
    if (!updateTodo.id || !updateTodo.todo.trim()) {
      setError('Both ID and todo text are required');
      return;
    }

    try {
      const response = await axios.put(`${apiUrl}/api/${backendName}/todo/${updateTodo.id}`, 
        { todo: updateTodo.todo }
      );
      
      if (response.data) {
        setTodos(
          todos.map((todo) => {
            if (todo.id === parseInt(updateTodo.id)) {
              return { ...todo, todo: updateTodo.todo };
            }
            return todo;
          })
        );
        setUpdateTodo({ id: '', todo: '' });
      }
    } catch (error: unknown) {
      if (error.response?.status === 404) {
        setError(`Todo with ID ${updateTodo.id} not found`);
      } else {
        setError(error.response?.data?.message || 'Failed to update todo');
      }
      console.error('Error updating todo:', error);
    }
  };

  const confirmDelete = (todoId: number) => {
    if (window.confirm('Are you sure you want to delete this todo?')) {
      deleteTodo(todoId);
    }
  };

  const deleteTodo = async (todoId: number) => {
    setError('');
    try {
      await axios.delete(`${apiUrl}/api/${backendName}/todo/${todoId}`);
      setTodos(todos.filter((todo) => todo.id !== todoId));
    } catch (error: unknown) {
      if (error.response?.status === 404) {
        setError(`Todo with ID ${todoId} not found`);
      } else {
        setError(error.response?.data?.message || 'Failed to delete todo');
      }
      console.error('Error deleting todo:', error);
    }
  };

  return (
    <div className={`todo-interface ${bgColor} ${backendName} w-full max-w-md p-4 my-4 rounded shadow`}>
      <h2 className="text-xl font-bold text-center text-white mb-6">{`Todo app`}</h2>

      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center p-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mx-auto"></div>
        </div>
      ) : (
        <>
          <form onSubmit={createTodo} className="mb-6 p-4 bg-blue-100 rounded shadow">
            <input
              placeholder="Add a new todo"
              value={newTodo.todo}
              onChange={(e) => setNewTodo({ todo: e.target.value })}
              className="mb-2 w-full p-2 border border-gray-300 rounded"
            />
            <button 
              type="submit" 
              className="w-full p-2 text-white bg-blue-500 rounded hover:bg-blue-600"
              disabled={!newTodo.todo.trim()}
            >
              Add Todo
            </button>
          </form>

          <form onSubmit={handleUpdateTodo} className="mb-6 p-4 bg-blue-100 rounded shadow">
            <input
              placeholder="Todo Id"
              value={updateTodo.id}
              onChange={(e) => setUpdateTodo({ ...updateTodo, id: e.target.value })}
              className="mb-2 w-full p-2 border border-gray-300 rounded"
              type="number"
            />
            <input
              placeholder="New Todo"
              value={updateTodo.todo}
              onChange={(e) => setUpdateTodo({ ...updateTodo, todo: e.target.value })}
              className="mb-2 w-full p-2 border border-gray-300 rounded"
            />
            <button 
              type="submit" 
              className="w-full p-2 text-white bg-green-500 rounded hover:bg-green-600"
              disabled={!updateTodo.id || !updateTodo.todo.trim()}
            >
              Update Todo
            </button>
          </form>


          <div className="space-y-4">
            {todos.length === 0 ? (
              <div className="text-center p-4 bg-gray-100 rounded">
                No todos found. Create one to get started!
              </div>
            ) : (
              todos.map((todo) => (
                <div key={todo.id} className="flex items-center justify-between bg-white p-4 rounded-lg shadow">
                  <CardComponent todo={todo} />
                  <button 
                    onClick={() => confirmDelete(todo.id)} 
                    className={`${btnColor} text-white py-2 px-4 rounded`}
                  >
                    Delete Todo
                  </button>
                </div>
              ))
            )}
          </div>
        </>
      )}
    </div>
  );
};

export default TodoInterface;
