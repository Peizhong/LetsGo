import { ADDTODO, TOGGLETODO, SETVISIBILITYFILTER } from "../constants"

let nextTodoId = 0

export interface IADDTODOAction {
    type: ADDTODO;
    id: number;
    text: string;
}

export interface ISETVISIBILITYFILTERAction{
    type: SETVISIBILITYFILTER;
    filter: string;
}

export interface ITOGGLEAction{
    type: TOGGLETODO;
    id: number;
}

export type ModifyAction = IADDTODOAction | ITOGGLEAction | ISETVISIBILITYFILTERAction;

export const addTodo = (text: string) => ({
  type: 'ADD_TODO',
  id: nextTodoId++,
  text
})

export const setVisibilityFilter = (filter: string) => ({
  type: 'SET_VISIBILITY_FILTER',
  filter
})

export const toggleTodo = (id: number) => ({
  type: 'TOGGLE_TODO',
  id
})

export const VisibilityFilters = {
  SHOW_ALL: 'SHOW_ALL',
  SHOW_COMPLETED: 'SHOW_COMPLETED',
  SHOW_ACTIVE: 'SHOW_ACTIVE'
}