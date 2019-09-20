import { ModifyAction } from "../actions"
import { ADDTODO, TOGGLETODO } from "../constants"
import { TODOSTATE } from "."

const initialTodoState: TODOSTATE = {
  todos: []
}

const todos = (state=initialTodoState, action: ModifyAction):TODOSTATE =>
{
  switch (action.type) {
    case ADDTODO:
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
    case TOGGLETODO:
      return {
        todos: state.todos.map((todo: any) => todo.id === action.id ? { ...todo, completed: !todo.completed } : todo)
      }
    default:
      return state
  }
}

export default todos