#include "parser.hpp"

#include <utility>
#include <iostream>

#include <boost/algorithm/string/trim.hpp>
#include <boost/algorithm/string/split.hpp>
#include <boost/range/algorithm_ext/erase.hpp>
#include <boost/algorithm/string/predicate.hpp>
#include <boost/program_options/variables_map.hpp>
#include <boost/algorithm/string/classification.hpp>

#include <util/read.hpp>
#include <parser/node/package_node.hpp>
#include <parser/node/def_function_node.hpp>
#include <parser/node/impl_function_node.hpp>

namespace ql::parser {
    Parser::Parser() {
        registerNode<PackageNode>("pckg");
        registerNode<ParseNode>("default");
        registerNode<DefineFunctionNode>("def");
        registerNode<ImplementFunctionNode>("impl");
    }

    std::shared_ptr<MasterNode> Parser::parse(po::variables_map& options) {
        auto sources = options["input"].as<std::vector<std::string>>();
        std::string sourceFileName = sources[0];
        auto src = util::readAllText(sourceFileName);
        auto node = getNodes(src.value());
        return node;
    }

    std::shared_ptr<AbstractNode> Parser::getNode
            (std::string const& nodeName, std::string&& blockWithInfo, std::vector<std::string>&& tokens, AbstractNode::ParentRef parent) {
        // Check if we have a generator function that can make this requested time, or else use default
        auto it = m_NamesToNodes.find(nodeName);
        nodeFactory& nodeFactoryFunc = it == m_NamesToNodes.end() ? m_NamesToNodes["default"] : it->second;
        auto node = nodeFactoryFunc(std::move(blockWithInfo), std::move(tokens), std::move(parent));
        return node;
    }

    std::shared_ptr<MasterNode> Parser::getNodes(std::string code) {
        auto parent = std::make_shared<MasterNode>();
        boost::remove_erase_if(code, boost::is_any_of("\n\r"));
        recurseNodes(code, parent);
        return parent;
    }

    void Parser::recurseNodes(std::string const& code, std::weak_ptr<AbstractNode> const& parent, int depth) {
        auto level = 0;
        int blockInfoStart = 0, blockStart = 0;
        for (int i = 0; i < static_cast<int>(code.size()); i++) {
            char c = code[i];
            if (c == '{') {
                if (level++ == 0) {
                    blockStart = i;
                }
            } else if (c == '}') {
                if (--level == 0) {
                    std::string blockWithInfo = code.substr(blockInfoStart, i - blockInfoStart + 1);
                    auto delimiters = boost::is_any_of("\t ");
                    // Trim is necessary since split will include empty strings in beginning if we do not
                    boost::trim_if(blockWithInfo, delimiters);
                    // Split by tabs and spaces into tokens, which we use to find what type of node to create
                    std::vector<std::string> tokens;
                    boost::split(tokens, blockWithInfo, delimiters, boost::token_compress_on);
                    std::string const& nodeName = tokens[0];
                    std::cout << nodeName << ": " << blockWithInfo << std::endl;
                    auto child = getNode(nodeName, std::move(blockWithInfo), std::move(tokens), parent);
                    // Add children to parent node, parent node is owning via a shared pointer
                    parent.lock()->addChild(child);
                    // Recurse on the inner contents of the block so that each node added is for one block only
                    std::string blockContents = code.substr(blockStart + 1, i - blockStart - 1);
                    recurseNodes(blockContents, child, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
