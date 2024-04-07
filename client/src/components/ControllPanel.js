import React, { useState } from "react";
import Box from '@mui/material/Box';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import axios from "axios";
import Paper from '@mui/material/Paper';

const ControllPanel = () => {
  const [operation, setOperation] = useState("PUT");
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [duration, setDuration] = useState("");

  const handleChange = (e) => {
    setOperation(e.target.value);
    setValue("");
    setKey("");
  }

  const handleSubmit = () => {
    if (operation == "PUT") {
      if (key == null || key == "" || duration == null || duration == "" || value == null || value == "") {
        return;
      }
      axios
        .put("http://localhost:8080/cache", {key: key, value: value, duration: duration})
        .then()
    } else if (operation == "GET") {
      if (key == null || key == "") {
        return;
      }
      axios
        .get(`http://localhost:8080/cache?key=${key}`)
        .then((res) => {
          if (res.status == 200) {
            alert(`received value for key '${key}': ${res.data.data}`);
          }
        })
    } else if (operation == "DELETE") {
      axios
        .delete(`http://localhost:8080/cache?key=${key}`)
        .then()
    }
  }

  return (
    <Box component={Paper}>
      <div style={{ width: "100%" }} className="card">
        <div className="header">
          Operations
        </div>
        <div style={{ width: "100%" }}>
          <Box sx={{ minWidth: 120, maxWidth: 200, margin: "auto", marginTop: "1rem" }}>
            <FormControl fullWidth>
              <InputLabel>Select Operation</InputLabel>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={operation}
                label="Select Operation"
                onChange={handleChange}
                size="small"
              >
                <MenuItem value={"PUT"}>PUT</MenuItem>
                <MenuItem value={"GET"}>GET</MenuItem>
                <MenuItem value={"DELETE"}>DELETE</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </div>
        <div style={{ width: "100%" }}>
          {
            operation == "PUT"
              ? PUTCMD()
              : operation == "GET"
                ? GETCMD()
                : DELETECMD()
          }
        </div>
        <div style={{ width: "100%" }}>
          <Button sx={{ marginTop: "1rem", marginBottom: "1rem" }} variant="contained" onClick={handleSubmit}>
            GO!
          </Button>
        </div>
      </div>
    </Box>
    
  );

  function PUTCMD() {
    return (
      <div style={{ display: "flex", flexDirection: "column", width: '100%' }}>
        <div style={{ width: "100%", display: "flex", flexDirection: "row", marginTop: "1rem", justifyContent: "space-evenly"  }}>
          <TextField label="key" variant="outlined" size="small" value={key} onChange={(e) => { setKey(e.target.value) }} />
          <TextField label="value" variant="outlined" size="small" value={value} onChange={(e) => { setValue(e.target.value) }} />
        </div>
        <div style={{ marginTop: "1rem" }}>
          <TextField label="duration" variant="outlined" size="small" value={duration} onChange={(e) => { setDuration(e.target.value) }} />
        </div>
      </div>
    );
  }

  function GETCMD() {
    return (
      <div style={{ display: "flex", flexDirection: "column", width: '100%' }}>
        <div style={{ width: "100%", display: "flex", flexDirection: "row", marginTop: "1rem", justifyContent: "space-evenly"  }}>
          <TextField label="key" variant="outlined" size="small" value={key} onChange={(e) => { setKey(e.target.value) }} />
        </div>
      </div>
    );
  }

  function DELETECMD() {
    return (
      <div style={{ display: "flex", flexDirection: "column", width: '100%' }}>
        <div style={{ width: "100%", display: "flex", flexDirection: "row", marginTop: "1rem", justifyContent: "space-evenly"  }}>
          <TextField label="key" variant="outlined" size="small" value={key} onChange={(e) => { setKey(e.target.value) }} />
        </div>
      </div>
    );
  }
}

export default ControllPanel;