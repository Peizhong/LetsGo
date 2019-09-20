import React from 'react'
import { TODO } from '../constants'
import Todo from './Todo'

interface TodoListProps{
  todos: TODO[]
}

// state和action
const TodoList = (p:TodoListProps) => (
    <ul>
      {p.todos.map(todo => (
      <Todo key={todo.id} {...todo} onClick={() => toggleTodo(todo.id)} />
      ))}
    </ul>
)