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
  type: ADDTODO,
  id: nextTodoId++,
  text
})

export const setVisibilityFilter = (filter: string) => ({
  type: SETVISIBILITYFILTER,
  filter
})

export const toggleTodo = (id: number) => ({
  type: TOGGLETODO,
  id
})

export const VisibilityFilters = {
  SHOW_ALL: 'SHOW_ALL',
  SHOW_COMPLETED: 'SHOW_COMPLETED',
  SHOW_ACTIVE: 'SHOW_ACTIVE'
}