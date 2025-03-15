import React, { useEffect, useState } from "react";
import {
  Container,
  TextField,
  Button,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Checkbox,
  Typography,
  Paper,
} from "@mui/material";
import { Delete as DeleteIcon } from "@mui/icons-material";

import axios from "axios";
import { toast, ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const API_URL = "http://localhost:4000/api/todo"; // Update if needed

const TodoApp = () => {
  const [todos, setTodos] = useState([]);
  const [newTodo, setNewTodo] = useState("");

  // Fetch Todos
  const fetchTodos = async () => {
    try {
      const response = await axios.get(API_URL);
      setTodos(response.data.todos);
    } catch (error) {
      console.error("Error fetching todos:", error);
      toast.error("Failed to fetch todos!");
    }
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  // Add New Todo
  const addTodo = async () => {
    if (!newTodo.trim()) {
      toast.warning("Please enter a todo!");
      return;
    }
    try {
      const response = await axios.post(API_URL, { body: newTodo });
      setTodos([...todos, response.data]);
      setNewTodo("");
      toast.success("Todo added!");
      fetchTodos();
    } catch (error) {
      console.error("Error adding todo:", error);
      toast.error("Failed to add todo!");
    }
  };

  // Mark Todo as Completed
  const completeTodo = async (id) => {
    try {
      await axios.patch(`${API_URL}/${id}`);
      setTodos(
        todos.map((todo) =>
          todo._id === id ? { ...todo, completed: true } : todo
        )
      );
      toast.info("Todo marked as completed!");
      fetchTodos();
    } catch (error) {
      console.error("Error updating todo:", error);
      toast.error("Failed to update todo!");
    }
  };

  // Delete Todo
  const deleteTodo = async (id) => {
    try {
      await axios.delete(`${API_URL}/${id}`);
      setTodos(todos.filter((todo) => todo._id !== id));
      toast.success("Todo deleted!");
      fetchTodos();
    } catch (error) {
      console.error("Error deleting todo:", error);
      toast.error("Failed to delete todo!");
    }
  };

  return (
    <Container maxWidth="sm">
      <ToastContainer />
      <Paper style={{ padding: 20, marginTop: 20 }}>
        <Typography variant="h4" align="center" gutterBottom>
          Todo App
        </Typography>
        <div style={{ display: "flex", gap: 10 }}>
          <TextField
            fullWidth
            variant="outlined"
            label="Add Todo"
            value={newTodo}
            onChange={(e) => setNewTodo(e.target.value)}
          />
          <Button variant="contained" color="primary" onClick={addTodo}>
            Add
          </Button>
        </div>
        <List>
          {todos.map((todo) => (
            <ListItem key={todo._id} divider>
              <Checkbox
                checked={todo.completed}
                onChange={() => completeTodo(todo.id)}
                disabled={todo.completed}
              />
              <ListItemText
                primary={todo.body}
                style={{
                  textDecoration: todo.completed ? "line-through" : "none",
                }}
              />
              <IconButton
                edge="end"
                color="error"
                onClick={() => {
                  console.log(todo);
                  deleteTodo(todo.id);
                }}
              >
                <DeleteIcon />
              </IconButton>
            </ListItem>
          ))}
        </List>
      </Paper>
    </Container>
  );
};

export default TodoApp;
