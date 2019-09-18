import { ModifyAction } from "../actions"
import { TODOSTATE } from "../constants"

const initialState: TODOSTATE = {
  todos: []
}

export default function todos(state=initialState, action: ModifyAction): TODOSTATE
{
  switch (action.type) {
    case 'ADD_TODO':
      return {
        todos:[
        ...state.todos,
        {
          id: action.id,
          text: action.text,
          completed: false,
          date: new Date()
        }]
      }
    case 'TOGGLE_TODO':
      return {
        todos: state.todos.map((todo: any) => todo.id === action.id ? { ...todo, completed: !todo.completed } : todo)
      }
    default:
      return state
  }
}