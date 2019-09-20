import { VisibilityFilters, ModifyAction } from '../actions'
import { SETVISIBILITYFILTER } from "../constants"
import { VIEWSTATE } from '.'

const initialViewState: VIEWSTATE = {
  filter: VisibilityFilters.SHOW_ALL
}

const visibilityFilter = (state = initialViewState, action: ModifyAction):VIEWSTATE => {
    switch (action.type) {
      case SETVISIBILITYFILTER:
        return {
          filter: action.filter
      }
      default:
        return state
    }
}

export default visibilityFilter