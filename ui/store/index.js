import _ from 'lodash';
import { EventEmitter } from 'events';

const store = new EventEmitter();

let alerts = [];

store.getAlerts = () => {
  return alerts;
};

store.deleteAlert = (id) => {
  alerts = _.filter(alerts, (alert) => {
    return alert.id !== id;
  });

  store.emit("alerts-changed");
};

store.createAlert = (msg, type) => {
  const id = _.uniqueId();
  alerts.splice(0, 0, { msg, type, id });
  window.setTimeout(() => store.deleteAlert(id), 6000);
  store.emit("alerts-changed");
};

export default store;
