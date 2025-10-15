// Workflow Builder Application
class WorkflowBuilder {
    constructor() {
        this.nodes = [];
        this.edges = [];
        this.selectedNode = null;
        this.nodeCounter = 0;
        this.connecting = null;
        
        this.canvas = document.getElementById('canvas');
        this.svg = document.getElementById('edgesSvg');
        this.propertiesPanel = document.getElementById('propertiesPanel');
        
        this.init();
    }
    
    init() {
        this.setupDragAndDrop();
        this.setupEventListeners();
        this.updateStats();
    }
    
    // Drag and Drop from Palette
    setupDragAndDrop() {
        const paletteNodes = document.querySelectorAll('.palette-node');
        
        paletteNodes.forEach(node => {
            node.addEventListener('dragstart', (e) => {
                e.dataTransfer.setData('nodeType', node.dataset.type);
            });
        });
        
        this.canvas.addEventListener('dragover', (e) => {
            e.preventDefault();
        });
        
        this.canvas.addEventListener('drop', (e) => {
            e.preventDefault();
            const nodeType = e.dataTransfer.getData('nodeType');
            const rect = this.canvas.getBoundingClientRect();
            const x = e.clientX - rect.left;
            const y = e.clientY - rect.top;
            
            this.addNode(nodeType, x, y);
        });
    }
    
    // Event Listeners
    setupEventListeners() {
        document.getElementById('clearBtn').addEventListener('click', () => this.clearAll());
        document.getElementById('loadBtn').addEventListener('click', () => this.openLoadModal());
        document.getElementById('saveBtn').addEventListener('click', () => this.saveWorkflow());
        document.getElementById('exportBtn').addEventListener('click', () => this.exportJSON());
        document.getElementById('executeBtn').addEventListener('click', () => this.executeWorkflow());
        document.getElementById('closePropertiesBtn').addEventListener('click', () => this.closeProperties());
        document.getElementById('deleteNodeBtn').addEventListener('click', () => this.deleteSelectedNode());
        document.getElementById('savePropertiesBtn').addEventListener('click', () => this.saveProperties());
        
        // File input handler
        document.getElementById('jsonFileInput').addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (file) {
                document.getElementById('fileName').textContent = file.name;
                const reader = new FileReader();
                reader.onload = (event) => {
                    document.getElementById('jsonInput').value = event.target.result;
                };
                reader.readAsText(file);
            }
        });
    }
    
    // Add Node to Canvas
    addNode(type, x, y) {
        const nodeId = `${type}-${++this.nodeCounter}`;
        const config = this.getDefaultConfig(type);
        
        const node = {
            id: nodeId,
            type: type,
            x: x,
            y: y,
            config: config
        };
        
        this.nodes.push(node);
        this.renderNode(node);
        this.updateStats();
        
        // Hide hint if first node
        if (this.nodes.length === 1) {
            const hint = document.querySelector('.canvas-hint');
            if (hint) hint.style.display = 'none';
        }
    }
    
    // Get Default Config for Node Type
    getDefaultConfig(type) {
        const configs = {
            start: {
                initialData: {
                    example: "value"
                }
            },
            condition: {
                lhs: "{{variable}}",
                operator: "==",
                rhs: "value"
            },
            mongodb_insert: {
                database: "workflow_db",
                collection: "collection_name",
                document: {
                    field: "{{value}}"
                }
            },
            mongodb_find: {
                database: "workflow_db",
                collection: "collection_name",
                filter: {
                    field: "{{value}}"
                },
                limit: 10,
                outputKey: "results"
            }
        };
        
        return configs[type] || {};
    }
    
    // Render Node on Canvas
    renderNode(node) {
        const nodeEl = document.createElement('div');
        nodeEl.className = 'canvas-node';
        nodeEl.setAttribute('data-node-id', node.id);
        nodeEl.style.left = node.x + 'px';
        nodeEl.style.top = node.y + 'px';
        
        const icons = {
            start: '▶️',
            condition: '❓',
            mongodb_insert: '💾',
            mongodb_find: '🔍'
        };
        
        const names = {
            start: 'Start',
            condition: 'Condition',
            mongodb_insert: 'MongoDB Insert',
            mongodb_find: 'MongoDB Find'
        };
        
        nodeEl.innerHTML = `
            <div class="connection-point input" data-node-id="${node.id}" data-point-type="input"></div>
            <div class="node-header">
                <span class="node-icon">${icons[node.type]}</span>
                <strong>${names[node.type]}</strong>
            </div>
            <div class="node-type">${node.type}</div>
            <div class="node-id">${node.id}</div>
            <div class="connection-point output" data-node-id="${node.id}" data-point-type="output"></div>
        `;
        
        // Make node draggable
        this.makeNodeDraggable(nodeEl, node);
        
        // Click to select and edit
        nodeEl.addEventListener('click', (e) => {
            if (!e.target.classList.contains('connection-point')) {
                this.selectNode(node);
            }
        });
        
        // Connection point events
        const connectionPoints = nodeEl.querySelectorAll('.connection-point');
        connectionPoints.forEach(point => {
            point.addEventListener('click', (e) => {
                e.stopPropagation();
                this.handleConnectionClick(point);
            });
        });
        
        this.canvas.appendChild(nodeEl);
    }
    
    // Make Node Draggable
    makeNodeDraggable(nodeEl, node) {
        let isDragging = false;
        let startX, startY, initialX, initialY;
        
        nodeEl.addEventListener('mousedown', (e) => {
            if (e.target.classList.contains('connection-point')) return;
            
            isDragging = true;
            startX = e.clientX;
            startY = e.clientY;
            initialX = node.x;
            initialY = node.y;
            nodeEl.classList.add('dragging');
        });
        
        document.addEventListener('mousemove', (e) => {
            if (!isDragging) return;
            
            const dx = e.clientX - startX;
            const dy = e.clientY - startY;
            
            node.x = initialX + dx;
            node.y = initialY + dy;
            
            nodeEl.style.left = node.x + 'px';
            nodeEl.style.top = node.y + 'px';
            
            this.redrawEdges();
        });
        
        document.addEventListener('mouseup', () => {
            if (isDragging) {
                isDragging = false;
                nodeEl.classList.remove('dragging');
            }
        });
    }
    
    // Handle Connection Point Click
    handleConnectionClick(point) {
        const nodeId = point.getAttribute('data-node-id');
        const type = point.getAttribute('data-point-type');
        
        if (!this.connecting) {
            // Start connection
            if (type === 'output') {
                this.connecting = { from: nodeId };
                point.style.background = '#48bb78';
                point.style.boxShadow = '0 0 0 4px rgba(72, 187, 120, 0.3)';
                console.log('Started connection from:', nodeId);
            }
        } else {
            // Complete connection
            if (type === 'input' && nodeId !== this.connecting.from) {
                this.connecting.to = nodeId;
                this.addEdge(this.connecting.from, this.connecting.to);
                
                // Reset
                const outputPoint = document.querySelector(`[data-node-id="${this.connecting.from}"][data-point-type="output"]`);
                if (outputPoint) {
                    outputPoint.style.background = '';
                    outputPoint.style.boxShadow = '';
                }
                this.connecting = null;
                console.log('Completed connection');
            } else {
                // Cancel connection
                const outputPoint = document.querySelector(`[data-node-id="${this.connecting.from}"][data-point-type="output"]`);
                if (outputPoint) {
                    outputPoint.style.background = '';
                    outputPoint.style.boxShadow = '';
                }
                this.connecting = null;
                console.log('Cancelled connection');
            }
        }
    }
    
    // Add Edge
    addEdge(fromNodeId, toNodeId) {
        // Determine output type based on node type
        const fromNode = this.nodes.find(n => n.id === fromNodeId);
        let output = 'default';
        
        if (fromNode && fromNode.type === 'condition') {
            // For condition, ask user
            output = prompt('Enter output type (true/false):', 'true') || 'true';
        }
        
        const edge = {
            from: fromNodeId,
            to: toNodeId,
            output: output
        };
        
        this.edges.push(edge);
        this.redrawEdges();
        this.updateStats();
    }
    
    // Redraw All Edges
    redrawEdges() {
        this.svg.innerHTML = '';
        
        this.edges.forEach((edge, index) => {
            const fromNode = this.nodes.find(n => n.id === edge.from);
            const toNode = this.nodes.find(n => n.id === edge.to);
            
            if (!fromNode || !toNode) return;
            
            const fromEl = document.querySelector(`[data-node-id="${edge.from}"]`);
            const toEl = document.querySelector(`[data-node-id="${edge.to}"]`);
            
            if (!fromEl || !toEl) return;
            
            const fromRect = fromEl.getBoundingClientRect();
            const toRect = toEl.getBoundingClientRect();
            const canvasRect = this.canvas.getBoundingClientRect();
            
            const x1 = fromNode.x + (fromRect.width) + 6;
            const y1 = fromNode.y + (fromRect.height / 2);
            const x2 = toNode.x - 6;
            const y2 = toNode.y + (toRect.height / 2);
            
            // Bezier curve
            const midX = (x1 + x2) / 2;
            const path = `M ${x1} ${y1} C ${midX} ${y1}, ${midX} ${y2}, ${x2} ${y2}`;
            
            const pathEl = document.createElementNS('http://www.w3.org/2000/svg', 'path');
            pathEl.setAttribute('d', path);
            pathEl.setAttribute('class', 'edge-line');
            pathEl.setAttribute('data-edge-index', index);
            
            // Add label
            const labelX = midX;
            const labelY = (y1 + y2) / 2;
            
            const label = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            label.setAttribute('x', labelX);
            label.setAttribute('y', labelY);
            label.setAttribute('class', 'edge-label');
            label.setAttribute('text-anchor', 'middle');
            label.textContent = edge.output;
            
            this.svg.appendChild(pathEl);
            this.svg.appendChild(label);
            
            // Click to delete edge
            pathEl.addEventListener('click', () => {
                if (confirm('Delete this connection?')) {
                    this.edges.splice(index, 1);
                    this.redrawEdges();
                    this.updateStats();
                }
            });
        });
    }
    
    // Select Node
    selectNode(node) {
        // Deselect all
        document.querySelectorAll('.canvas-node').forEach(n => n.classList.remove('selected'));
        
        // Select this one
        const nodeEl = document.querySelector(`[data-node-id="${node.id}"]`);
        if (nodeEl) nodeEl.classList.add('selected');
        
        this.selectedNode = node;
        this.showProperties(node);
    }
    
    // Show Properties Panel
    showProperties(node) {
        this.propertiesPanel.classList.add('active');
        
        const content = document.getElementById('propertiesContent');
        content.innerHTML = this.getPropertiesForm(node);
    }
    
    // Get Properties Form HTML
    getPropertiesForm(node) {
        let formHTML = `
            <div class="form-group">
                <label>Node ID</label>
                <input type="text" id="prop-id" value="${node.id}" readonly>
            </div>
            <div class="form-group">
                <label>Node Type</label>
                <input type="text" value="${node.type}" readonly>
            </div>
        `;
        
        if (node.type === 'start') {
            formHTML += `
                <div class="form-group">
                    <label>Initial Data (JSON)</label>
                    <textarea id="prop-config">${JSON.stringify(node.config.initialData, null, 2)}</textarea>
                </div>
            `;
        } else if (node.type === 'condition') {
            formHTML += `
                <div class="form-group">
                    <label>Left Side (LHS)</label>
                    <input type="text" id="prop-lhs" value="${node.config.lhs}">
                </div>
                <div class="form-group">
                    <label>Operator</label>
                    <select id="prop-operator">
                        <option value="==" ${node.config.operator === '==' ? 'selected' : ''}>==</option>
                        <option value="!=" ${node.config.operator === '!=' ? 'selected' : ''}>!=</option>
                        <option value=">" ${node.config.operator === '>' ? 'selected' : ''}>&gt;</option>
                        <option value="<" ${node.config.operator === '<' ? 'selected' : ''}>&lt;</option>
                        <option value=">=" ${node.config.operator === '>=' ? 'selected' : ''}>&gt;=</option>
                        <option value="<=" ${node.config.operator === '<=' ? 'selected' : ''}>&lt;=</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Right Side (RHS)</label>
                    <input type="text" id="prop-rhs" value="${node.config.rhs}">
                </div>
            `;
        } else if (node.type === 'mongodb_insert') {
            formHTML += `
                <div class="form-group">
                    <label>Database</label>
                    <input type="text" id="prop-database" value="${node.config.database}">
                </div>
                <div class="form-group">
                    <label>Collection</label>
                    <input type="text" id="prop-collection" value="${node.config.collection}">
                </div>
                <div class="form-group">
                    <label>Document (JSON)</label>
                    <textarea id="prop-document">${JSON.stringify(node.config.document, null, 2)}</textarea>
                </div>
            `;
        } else if (node.type === 'mongodb_find') {
            formHTML += `
                <div class="form-group">
                    <label>Database</label>
                    <input type="text" id="prop-database" value="${node.config.database}">
                </div>
                <div class="form-group">
                    <label>Collection</label>
                    <input type="text" id="prop-collection" value="${node.config.collection}">
                </div>
                <div class="form-group">
                    <label>Filter (JSON)</label>
                    <textarea id="prop-filter">${JSON.stringify(node.config.filter, null, 2)}</textarea>
                </div>
                <div class="form-group">
                    <label>Limit</label>
                    <input type="number" id="prop-limit" value="${node.config.limit || 10}">
                </div>
                <div class="form-group">
                    <label>Output Key</label>
                    <input type="text" id="prop-outputKey" value="${node.config.outputKey || 'results'}">
                </div>
            `;
        }
        
        return formHTML;
    }
    
    // Save Properties
    saveProperties() {
        if (!this.selectedNode) return;
        
        const node = this.selectedNode;
        
        if (node.type === 'start') {
            try {
                const data = JSON.parse(document.getElementById('prop-config').value);
                node.config.initialData = data;
            } catch (e) {
                alert('Invalid JSON format');
                return;
            }
        } else if (node.type === 'condition') {
            node.config.lhs = document.getElementById('prop-lhs').value;
            node.config.operator = document.getElementById('prop-operator').value;
            node.config.rhs = document.getElementById('prop-rhs').value;
        } else if (node.type === 'mongodb_insert') {
            node.config.database = document.getElementById('prop-database').value;
            node.config.collection = document.getElementById('prop-collection').value;
            try {
                node.config.document = JSON.parse(document.getElementById('prop-document').value);
            } catch (e) {
                alert('Invalid JSON in document field');
                return;
            }
        } else if (node.type === 'mongodb_find') {
            node.config.database = document.getElementById('prop-database').value;
            node.config.collection = document.getElementById('prop-collection').value;
            try {
                node.config.filter = JSON.parse(document.getElementById('prop-filter').value);
            } catch (e) {
                alert('Invalid JSON in filter field');
                return;
            }
            node.config.limit = parseInt(document.getElementById('prop-limit').value);
            node.config.outputKey = document.getElementById('prop-outputKey').value;
        }
        
        alert('Properties saved!');
    }
    
    // Close Properties Panel
    closeProperties() {
        this.propertiesPanel.classList.remove('active');
        this.selectedNode = null;
        document.querySelectorAll('.canvas-node').forEach(n => n.classList.remove('selected'));
    }
    
    // Delete Selected Node
    deleteSelectedNode() {
        if (!this.selectedNode) return;
        
        if (!confirm(`Delete node ${this.selectedNode.id}?`)) return;
        
        // Remove node
        const nodeId = this.selectedNode.id;
        this.nodes = this.nodes.filter(n => n.id !== nodeId);
        
        // Remove connected edges
        this.edges = this.edges.filter(e => e.from !== nodeId && e.to !== nodeId);
        
        // Remove from DOM
        const nodeEl = document.querySelector(`[data-node-id="${nodeId}"]`);
        if (nodeEl) nodeEl.remove();
        
        this.closeProperties();
        this.redrawEdges();
        this.updateStats();
    }
    
    // Update Stats
    updateStats() {
        document.getElementById('nodeCount').textContent = this.nodes.length;
        document.getElementById('edgeCount').textContent = this.edges.length;
    }
    
    // Clear All
    clearAll() {
        if (!confirm('Clear all nodes and edges?')) return;
        
        this.nodes = [];
        this.edges = [];
        this.selectedNode = null;
        this.nodeCounter = 0;
        
        // Clear canvas
        this.canvas.querySelectorAll('.canvas-node').forEach(n => n.remove());
        this.svg.innerHTML = '';
        
        // Show hint again
        const hint = document.querySelector('.canvas-hint');
        if (hint) hint.style.display = 'block';
        
        this.closeProperties();
        this.updateStats();
    }
    
    // Export JSON
    exportJSON() {
        const workflowId = document.getElementById('workflowId').value;
        const workflowName = document.getElementById('workflowName').value;
        
        const workflow = {
            id: workflowId,
            name: workflowName,
            nodes: this.nodes.map(n => ({
                id: n.id,
                type: n.type,
                config: n.config
            })),
            edges: this.edges
        };
        
        const json = JSON.stringify(workflow, null, 2);
        document.getElementById('jsonOutput').value = json;
        document.getElementById('jsonModal').classList.add('active');
    }
    
    // Save Workflow to Database
    async saveWorkflow() {
        if (this.nodes.length === 0) {
            alert('Please add at least one node before saving!');
            return;
        }
        
        const workflowId = document.getElementById('workflowId').value;
        const workflowName = document.getElementById('workflowName').value;
        
        if (!workflowName.trim()) {
            alert('Please enter a workflow name!');
            return;
        }
        
        const workflow = {
            id: workflowId,
            name: workflowName,
            nodes: this.nodes.map(n => ({
                id: n.id,
                type: n.type,
                config: n.config
            })),
            edges: this.edges
        };
        
        try {
            const response = await fetch('http://localhost:3002/create-workflow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(workflow)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                // Update workflow ID with the one from server
                if (result.workflow_id) {
                    document.getElementById('workflowId').value = result.workflow_id;
                }
                alert(`✅ Workflow saved successfully!\nID: ${result.workflow_id}`);
            } else {
                alert('Error saving workflow: ' + (result.details || result.error));
            }
        } catch (error) {
            alert('Error saving workflow: ' + error.message);
        }
    }
    
    // Execute Workflow
    async executeWorkflow() {
        if (this.nodes.length === 0) {
            alert('Please add at least one node before executing!');
            return;
        }
        
        const workflowId = document.getElementById('workflowId').value;
        const workflowName = document.getElementById('workflowName').value;
        
        const workflow = {
            id: workflowId,
            name: workflowName,
            nodes: this.nodes.map(n => ({
                id: n.id,
                type: n.type,
                config: n.config
            })),
            edges: this.edges
        };
        
        try {
            const response = await fetch('http://localhost:3002/execute-workflow', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(workflow)
            });
            
            const result = await response.json();
            
            document.getElementById('resultOutput').textContent = JSON.stringify(result, null, 2);
            document.getElementById('resultModal').classList.add('active');
        } catch (error) {
            alert('Error executing workflow: ' + error.message);
            document.getElementById('resultOutput').textContent = 'Error: ' + error.message;
            document.getElementById('resultModal').classList.add('active');
        }
    }
    
    // Open Load Modal
    openLoadModal() {
        document.getElementById('jsonInput').value = '';
        document.getElementById('fileName').textContent = '';
        document.getElementById('jsonFileInput').value = '';
        document.getElementById('loadModal').classList.add('active');
    }
    
    // Load Workflow from JSON
    loadFromJSON(workflow) {
        // Clear current workflow
        this.nodes = [];
        this.edges = [];
        this.canvas.querySelectorAll('.canvas-node').forEach(n => n.remove());
        this.svg.innerHTML = '';
        
        // Set workflow info
        document.getElementById('workflowId').value = workflow.id || 'workflow-1';
        document.getElementById('workflowName').value = workflow.name || 'My Workflow';
        
        // Load nodes
        workflow.nodes.forEach((nodeDef, index) => {
            // Calculate position in a grid layout
            const x = 100 + (index % 3) * 250;
            const y = 100 + Math.floor(index / 3) * 150;
            
            const node = {
                id: nodeDef.id,
                type: nodeDef.type,
                x: x,
                y: y,
                config: nodeDef.config || this.getDefaultConfig(nodeDef.type)
            };
            
            this.nodes.push(node);
            this.renderNode(node);
        });
        
        // Load edges
        workflow.edges.forEach(edgeDef => {
            this.edges.push({
                from: edgeDef.from,
                to: edgeDef.to,
                output: edgeDef.output || 'default'
            });
        });
        
        // Update counter to avoid ID conflicts
        this.nodeCounter = this.nodes.length;
        
        // Redraw edges
        setTimeout(() => {
            this.redrawEdges();
            this.updateStats();
        }, 100);
        
        // Hide hint
        const hint = document.querySelector('.canvas-hint');
        if (hint) hint.style.display = 'none';
        
        console.log('Loaded workflow:', workflow.name);
    }
}

// Modal Functions
function closeModal() {
    document.getElementById('jsonModal').classList.remove('active');
}

function closeResultModal() {
    document.getElementById('resultModal').classList.remove('active');
}

function closeLoadModal() {
    document.getElementById('loadModal').classList.remove('active');
}

function copyJSON() {
    const textarea = document.getElementById('jsonOutput');
    textarea.select();
    document.execCommand('copy');
    alert('JSON copied to clipboard!');
}

function downloadJSON() {
    const json = document.getElementById('jsonOutput').value;
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'workflow.json';
    a.click();
    URL.revokeObjectURL(url);
}

function loadWorkflowJSON() {
    const jsonInput = document.getElementById('jsonInput').value.trim();
    
    if (!jsonInput) {
        alert('Please paste JSON or select a file!');
        return;
    }
    
    try {
        const workflow = JSON.parse(jsonInput);
        
        // Validate workflow structure
        if (!workflow.nodes || !Array.isArray(workflow.nodes)) {
            alert('Invalid workflow JSON: missing or invalid "nodes" array');
            return;
        }
        
        if (!workflow.edges || !Array.isArray(workflow.edges)) {
            alert('Invalid workflow JSON: missing or invalid "edges" array');
            return;
        }
        
        // Validate nodes
        for (const node of workflow.nodes) {
            if (!node.id || !node.type) {
                alert('Invalid node: missing id or type');
                return;
            }
            
            const validTypes = ['start', 'condition', 'mongodb_insert', 'mongodb_find'];
            if (!validTypes.includes(node.type)) {
                alert(`Invalid node type: ${node.type}. Valid types are: ${validTypes.join(', ')}`);
                return;
            }
        }
        
        // Load the workflow
        app.loadFromJSON(workflow);
        closeLoadModal();
        alert('✅ Workflow loaded successfully!');
        
    } catch (error) {
        alert('Error parsing JSON: ' + error.message + '\n\nPlease check your JSON format.');
        console.error('JSON parse error:', error);
    }
}

// Initialize App
const app = new WorkflowBuilder();

