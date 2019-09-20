export const ADDTODO = "ADD_TODO";
export type ADDTODO = typeof ADDTODO;

export const SETVISIBILITYFILTER = "SET_VISIBILITY_FILTER"
export type SETVISIBILITYFILTER = typeof SETVISIBILITYFILTER;

export const TOGGLETODO = "TOGGLE_TODO";
export type TOGGLETODO = typeof TOGGLETODO;

export interface TODO{
    id: number;
    text: string;
    completed: boolean;
    date: Date;
}