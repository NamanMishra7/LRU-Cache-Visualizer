import React, { useEffect, useState } from "react";
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

const Cache = () => {
  const [socket, setSocket] = useState(null);
  const [cache, setCache] = useState([]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    setSocket(ws);

    ws.onmessage = event => {
      const message = JSON.parse(event.data);
      setCache(JSON.parse(message.content));
      // console.log(JSON.parse(message.content))
      // setMessages(prevMessages => [...prevMessages, message]);
    };

    // return () => {
    //   ws.close();
    // };
  }, []);



  return (
    <div style={{ width: "100%", fontSize: "120%" }}>
      <TableContainer component={Paper} sx={{ maxWidth: 650, margin: "auto" }}>
      <Table sx={{ maxWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell><b>Key</b></TableCell>
            <TableCell align="right"><b>Value</b></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {cache.map((row, index) => (
            <TableRow
              key={index}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell component="th" scope="row">
                {Object.keys(row)[0]}
              </TableCell>
              <TableCell align="right">{Object.values(row)[0]}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
    </div>
  );
}

export default Cache;