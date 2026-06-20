import {info} from "../site-wide/util";

export default class Publisher {

    constructor() {
        this.lastEvents = {};
        this.subs = {};
    }

    /**
     *
     * @param chan
     * @param targetId (any|null) null will broadcast
     * @param val
     * @param debounce
     */
    notify(chan, targetId, val, debounce) {
        info(chan, targetId, val);
        if (debounce && this.lastEvents[chan] === val) return;
        for (let id in this.subs[chan]) {
            if (!this.subs[chan].hasOwnProperty(id)) continue;
            if (targetId !== null) {
                if (targetId === id) {
                    this.subs[chan][id](chan, val);
                }
                continue;
            }
            this.subs[chan][id](chan, val);
        }
        this.lastEvents[chan] = val;
    }

    /**
     *
     * @param channels
     * @param id identifier of subscriber
     * @param fn(chan, val) function to call
     */
    subscribe(channels, id, fn) {
        channels.forEach(c => {
            this.subs[c] = this.subs[c] || {};
            this.subs[c][id] = fn;
        });
    }

}