import { combineReducers } from 'redux'
import todos from './todos'
import visibilityFilter from './visibilityFilter'
import { TODO } from '../constants'

export default combineReducers({
    todos,
    visibilityFilter
})

export interface TODOSTATE{
    todos: TODO[];
}

export interface VIEWSTATE{
    filter: string;
}