import http from "k6/http";
import { sleep } from 'k6';

export const options = {
  iterations: 1,
};

export default function () {
  const r = {
    method: 'GET',
    url: 'http://localhost:1729',
  };
  const responses = http.batch([r, r, r]);
  sleep(1);
}
