
export  function handleJsonFromString(str) {
    let startIndexObject = str.indexOf('{');
    let startIndexArray = str.indexOf('[');
    let jsonStartIndex = -1;
  
    if (startIndexObject === -1 && startIndexArray === -1) {

      console.error("No JSON start characters '{' or '[' found.");
      return null;
    } else if (startIndexObject === -1) {
      jsonStartIndex = startIndexArray;
    } else if (startIndexArray === -1) {
      jsonStartIndex = startIndexObject;
    } else {
      jsonStartIndex = Math.min(startIndexObject, startIndexArray);
    }
  
    let endIndexObject = str.lastIndexOf('}');
    let endIndexArray = str.lastIndexOf(']');
    let jsonEndIndex = -1;
  

    if (endIndexObject === -1 && endIndexArray === -1) {

      console.error("No JSON end characters '}' or ']' found.");
      return null;
    } else if (endIndexObject === -1) {
      jsonEndIndex = endIndexArray;
    } else if (endIndexArray === -1) {
      jsonEndIndex = endIndexObject;
    } else {
      jsonEndIndex = Math.max(endIndexObject, endIndexArray);
    }
  
    if (jsonStartIndex === -1 || jsonEndIndex === -1 || jsonEndIndex < jsonStartIndex) {
      console.error("Could not determine a valid JSON start/end range.");
      return null;
    }
  
    // Extract the potential JSON string
    const potentialJson = str.substring(jsonStartIndex, jsonEndIndex + 1);
  
    try {
      // Attempt to parse it
      JSON.parse(potentialJson);

      return potentialJson;
    } catch (e) {
      console.error("Extracted string is not valid JSON:", e.message);
      console.error("Extracted string was:", potentialJson);

      return null;
    }
  }