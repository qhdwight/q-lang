#include "parser.hpp"

#include <utility>
#include <iostream>

#include <boost/algorithm/string/split.hpp>
#include <boost/algorithm/string/trim_all.hpp>
#include <boost/range/algorithm_ext/erase.hpp>
#include <boost/algorithm/string/predicate.hpp>
#include <boost/program_options/variables_map.hpp>
#include <boost/algorithm/string/classification.hpp>

#include <util/read.hpp>
#include <parser/node/structure/package_node.hpp>
#include <parser/node/structure/def_func_node.hpp>
#include <parser/node/structure/impl_func_node.hpp>

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
        std::cout << sourceFileName << std::endl;
        auto src = util::readAllText(sourceFileName);
        auto node = getNodes(src.value());
        return node;
    }

    std::shared_ptr<AbstractNode> Parser::getNode(std::string const& nodeName,
                                                  std::string&& blockWithInfo, std::string_view const& innerBlock, std::vector<std::string>&& tokens,
                                                  AbstractNode::ParentRef parent) {
        // Check if we have a generator function that can make this requested time, or else use default
        auto it = m_NamesToNodes.find(nodeName);
        NodeFactory& nodeFactoryFunc = it == m_NamesToNodes.end() ? m_NamesToNodes["default"] : it->second;
        auto node = nodeFactoryFunc(std::move(blockWithInfo), innerBlock, std::move(tokens), std::move(parent));
        node->parse();
        return node;
    }

    std::shared_ptr<MasterNode> Parser::getNodes(std::string code) {
        auto parent = std::make_shared<MasterNode>();
        boost::remove_erase_if(code, boost::is_any_of("\n\r"));
        recurseNodes(code, parent);
        return parent;
    }

    void Parser::recurseNodes(std::string_view code, std::weak_ptr<AbstractNode> const& parent, int depth) {
        auto level = 0;
        auto blockInfoStart = 0ul, blockStart = 0ul;
        for (auto i = 0ul; i < code.size(); i++) {
            char c = code[i];
            if (c == '{') {
                if (level++ == 0) {
                    blockStart = i + 1ul;
                }
            } else if (c == '}') {
                if (--level == 0) {
                    auto blockInfoSize = i - blockInfoStart + 1ul;
                    std::string blockWithInfo(code.substr(blockInfoStart, blockInfoSize));

                    // Split by tabs and spaces into tokens, which we use to find what type of node to create
                    std::vector<std::string> tokens;
                    auto deliminator = boost::is_any_of("\t ");
                    boost::split(tokens, blockWithInfo, deliminator, boost::token_compress_on);
                    // Remove first and last blank tokens if they exist
                    boost::trim_all_if(tokens, [](auto const& token) { return token.empty(); });

                    std::string const& nodeName = tokens[0ul]; // TODO do more checks as opposed to just taking first
                    std::cout << nodeName << " --> " << blockWithInfo << "#" << std::endl;

                    // Find inner block
                    auto blockContentStart = blockStart, blockContentSize = i - blockContentStart;
                    std::string_view blockContents = std::string_view(blockWithInfo).substr(blockContentStart - blockInfoStart, blockContentSize);
                    std::cout << blockContents << "#" << std::endl;
                    auto child = getNode(nodeName, std::move(blockWithInfo), blockContents, std::move(tokens), parent);
                    // Add children to parent node, parent node is owning via a shared pointer
                    parent.lock()->addChild(child);

                    // Recurse on the inner contents of the block so that each node added is for one block only
                    recurseNodes(code.substr(blockContentStart, blockContentSize), child, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
