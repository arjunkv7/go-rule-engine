# Workflow Builder UI

A visual workflow builder for the Go Workflow Engine.

## Features

- 🎨 **Drag & Drop Interface** - Easily add nodes to the canvas
- 🔗 **Visual Connections** - Click output/input points to connect nodes
- ⚙️ **Node Configuration** - Edit node properties in the side panel
- 📋 **JSON Export** - Generate workflow JSON
- ▶️ **Direct Execution** - Test workflows directly from the UI
- 💾 **Download** - Save workflows as JSON files

## How to Use

### 1. Start the Server

```bash
cd /home/arjun/Projects/go-workflow-engine
go run main.go handlers.go
```

### 2. Open the UI

Navigate to: http://localhost:3002/ui

### 3. Build Your Workflow

#### Adding Nodes
1. Drag a node from the left sidebar
2. Drop it on the canvas
3. The node appears with connection points

#### Connecting Nodes
1. Click the **output point** (● on right side) of a source node
   - The point will turn green when active
2. Click the **input point** (● on left side) of a target node
3. For condition nodes, you'll be prompted to specify the output type (true/false)
4. The connection will appear as a curved line
5. Click on any edge to delete it

#### Configuring Nodes
1. Click on any node to select it
2. The properties panel opens on the right
3. Edit the configuration
4. Click **Save** to apply changes

#### Node Types

**Start Node** ▶️
- Entry point of the workflow
- Configure initial data in JSON format

**Condition Node** ❓
- Conditional branching
- Configure: LHS, Operator, RHS
- Outputs: "true" or "false"

**MongoDB Insert** 💾
- Insert documents into MongoDB
- Configure: database, collection, document
- Use `{{variableName}}` for dynamic values

**MongoDB Find** 🔍
- Query documents from MongoDB
- Configure: database, collection, filter, limit, outputKey
- Results stored in context with the specified outputKey

### 4. Managing Workflows

#### Load Existing Workflow
1. Click **📂 Load JSON** button
2. Either:
   - **Paste JSON** directly into the textarea, or
   - **Click "Choose File"** to upload a .json file
3. Click **Load Workflow**
4. The workflow appears on the canvas with all nodes and connections
5. You can now edit, save, or execute it

**Supports:**
- Loading from exported JSON files
- Loading from example workflows
- Importing workflows from other sources
- Validation of JSON structure before loading

#### Save Workflow to Database
1. Enter a workflow name
2. Click **💾 Save to DB** button
3. The workflow is saved to MongoDB
4. A unique ID is generated and displayed
5. Use this ID to execute the workflow later via API

#### Export JSON
1. Click **Export JSON** button
2. View the generated JSON
3. Copy to clipboard or download
4. Use for version control or sharing

#### Execute Workflow
1. Click **Execute Workflow** button
2. The workflow runs on the server
3. View the result in the modal
4. Check for any errors in the response

#### Clear All
1. Click **Clear All** to reset the canvas
2. Confirms before deleting all nodes and edges

### 5. Workflow Information

- **Workflow ID**: Unique identifier for the workflow
- **Workflow Name**: Human-readable name
- **Stats**: Shows node and edge counts

## Tips & Tricks

### Template Variables

Use `{{variableName}}` to access data from the workflow context:

```json
{
  "name": "{{userName}}",
  "age": "{{userAge}}"
}
```

### Multiple Connections

- One node can have multiple outgoing connections
- Connections with the same output execute in parallel

### Deleting Elements

- **Delete Node**: Select node → Click "Delete Node" in properties panel
- **Delete Edge**: Click on the edge line → Confirm deletion

### Keyboard Shortcuts

- **Escape**: Close properties panel
- **Delete**: Delete selected node (when properties panel is open)

## Example Workflows

### Simple User Registration

1. Add **Start** node with initial user data
2. Add **Condition** node to check age >= 18
3. Add **MongoDB Insert** to save user if adult
4. Connect: Start → Condition (default) → Insert (true)

### Find and Report

1. Add **Start** node with search criteria
2. Add **MongoDB Find** to query users
3. Add **Condition** to check if results found
4. Add **MongoDB Insert** to create report
5. Connect: Start → Find → Condition → Insert

## Troubleshooting

### Nodes won't connect
- **Solution**: Click the **output point** (right side, ●) first to start connection
  - It will turn green when active
  - Then click the **input point** (left side, ●) of another node
  - Check browser console (F12) for connection logs
- **Note**: Can't connect a node to itself
- **Tip**: Hover over connection points - they should grow and turn green

### Edge/Connection not visible
- The edge might be outside the canvas view
- Try moving nodes closer together
- Check that both nodes exist before connecting
- Look for the edge in the SVG layer (it may be very light gray)

### Workflow won't save to DB
- Ensure at least one node is added to the canvas
- Enter a workflow name before saving
- Check that MongoDB is running
- Verify the create-workflow API endpoint is working

### Workflow won't execute
- Check that MongoDB is running
- Verify API server is running on port 3002
- Check browser console (F12) for errors
- Ensure Start node is connected to other nodes
- Verify all required node configurations are filled

### JSON validation errors
- Ensure JSON fields in node properties are valid
- Check for missing quotes or commas
- Use a JSON validator if needed
- Try the "Export JSON" button to see the generated workflow

## File Structure

```
ui/
├── index.html       # Main UI page
├── css/
│   └── style.css   # Styling
└── js/
    └── app.js      # Workflow builder logic
```

## Browser Support

- Chrome (recommended)
- Firefox
- Safari
- Edge

## Future Enhancements

- [ ] Zoom and pan canvas
- [ ] Undo/redo functionality
- [ ] Workflow templates
- [ ] Auto-save to localStorage
- [ ] Keyboard shortcuts
- [ ] Multi-select nodes
- [ ] Copy/paste nodes
- [ ] Grid snapping
- [ ] Mini-map
- [ ] Search nodes

---

**Enjoy building workflows!** 🚀

