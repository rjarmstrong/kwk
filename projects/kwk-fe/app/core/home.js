import {decodeToken} from '../api/client'
import PubRoot from './pubRoot'
import {info} from '../site-wide/util';


export default class Home extends Component {

        constructor(props) {
            super(props);
            info()
        }
        render() {
            info();
            let t = decodeToken();
            if (t.key==='anon') {
                return <PubRoot/>
            }
            this.context.router.push('/' + t.key)
        }
}