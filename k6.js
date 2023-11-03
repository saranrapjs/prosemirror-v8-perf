import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.post('http://localhost:8080/', `{"steps":[{"stepType":"replace","from":0,"to":0,"slice":{"content":[{"type":"text","text":"Hello"}]}}],"doc": {"type":"doc"}}`, {
     headers: { 'Content-Type': 'application/json' },
  });
  sleep(1);
}
